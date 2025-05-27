package server

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/core/blockchain"
	"github.com/andantan/vote-blockchain-server/core/mempool"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	EXIT_SIGNAL = iota
	CONVERT_SIGNAL
	STANDBY_SIGNAL
)

const (
	GRPC_REQUEST_BUFFER_SIZE = 1024
)

// gRPC Network and port options
//
// Network must be "tcp", "unix" or "unixpacket"
//
// Port must be between 0 and 65535
type BlockChainListenerOpts struct {
	TopicgRPCNetwork     string
	TopicgRPCNetworkPort uint16
	VotegRPCNetwork      string
	VotegRPCNetworkPort  uint16
}

func (o *BlockChainServerOpts) SetTopicOptions(topicNetwork string, topicNetworkPort uint16) {
	o.TopicgRPCNetwork = topicNetwork
	o.TopicgRPCNetworkPort = topicNetworkPort
}

func (o *BlockChainServerOpts) SetVoteOptions(voteNetwork string, voteNetworkPort uint16) {
	o.VotegRPCNetwork = voteNetwork
	o.VotegRPCNetworkPort = voteNetworkPort
}

type BlockChainControllOpts struct {
	BlockTime time.Duration
	MaxTxSize uint32
}

func (o *BlockChainServerOpts) SetControllOptions(blockTime time.Duration, maxTxSize uint32) {
	o.BlockTime = blockTime
	o.MaxTxSize = maxTxSize
}

type BlockChainServerOpts struct {
	BlockChainListenerOpts
	BlockChainControllOpts
}

func NewBlockChainServerOpts() BlockChainServerOpts {
	return BlockChainServerOpts{}
}

type BlockChainServer struct {
	BlockChainServerOpts

	*BlockChainVoteServer
	*BlockChainTopicServer

	mempool    *mempool.MemPool
	blockChain *blockchain.BlockChain

	pendedCh     <-chan *mempool.Pended
	protoBlockCh chan<- *block.ProtoBlock
	ExitSignalCh chan uint8
}

func NewBlockChainServer(opts BlockChainServerOpts) *BlockChainServer {
	server := &BlockChainServer{
		BlockChainServerOpts: opts,
	}

	server.Initialize()

	log.Println(util.SystemString("SYSTEM: BlockChainServer engine generated"))

	return server
}

func (s *BlockChainServer) Initialize() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer initialize..."))

	s.setgRPCServer()
	s.setMemPool()
	s.setBlockChain()
	s.setChannel()

	log.Println(util.SystemString("SYSTEM: BlockChainServer initialization is done."))
}

func (s *BlockChainServer) setgRPCServer() {
	log.Printf(
		util.SystemString("SYSTEM: BlockChainServer setting gRPC server... | { GRPC_REQUEST_BUFFER_SIZE: %d }"),
		GRPC_REQUEST_BUFFER_SIZE,
	)

	s.BlockChainVoteServer = NewBlockChainVoteServer(GRPC_REQUEST_BUFFER_SIZE)
	s.BlockChainTopicServer = NewBlockChainTopicServer(GRPC_REQUEST_BUFFER_SIZE)

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting gRPC server is done."))
}

func (s *BlockChainServer) setMemPool() {
	log.Printf(
		util.SystemString("SYSTEM: BlockChainServer setting memory pool... | { BlockTime: %s, MaxTxSize: %d }"),
		s.BlockChainControllOpts.BlockTime,
		s.BlockChainControllOpts.MaxTxSize,
	)

	s.mempool = mempool.NewMemPool(
		s.BlockChainControllOpts.BlockTime,
		s.BlockChainControllOpts.MaxTxSize,
	)

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting memory pool is done."))
}

func (s *BlockChainServer) setBlockChain() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting blockchain..."))

	s.blockChain = blockchain.NewBlockChainWithGenesisBlock()

	genesisHeader, err := s.blockChain.GetHeader(0)

	if err != nil {
		log.Fatalf(util.RedString("Genesis block initialization error: %s"), err.Error())
	}

	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block ID=%s"), genesisHeader.VotingID)
	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block MerkleRoot=%s"), genesisHeader.MerkleRoot.String())
	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block Height=%d"), genesisHeader.Height)
	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block PrevBlockHash=%s"), genesisHeader.PrevBlockHash.String())
	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block BlockHash=%s"), genesisHeader.Hash().String())
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting blockchain is done."))
}

func (s *BlockChainServer) setChannel() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting channel..."))

	s.mempool.SetChannel()

	s.pendedCh = s.mempool.Produce()
	s.protoBlockCh = s.blockChain.ProtoBlockProducer()

	s.ExitSignalCh = make(chan uint8)

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting channel is done."))
}

func (s *BlockChainServer) Start() {
	s.startgRPCListener()

labelServer:
	for {
		select {
		case topic := <-s.RequestTopicCh:

			// TODO Make this tender
			if strings.Compare(string(topic.Topic), "exit") == 0 {
				topic.ResponseCh <- s.GetSuccessSubmitTopic("shutdown")

				s.BlockChainTopicServer.Shutdown()
				s.BlockChainVoteServer.Shutdown()
				s.mempool.Shutdown()
				s.blockChain.Shutdown()

				break labelServer
			}

			if err := s.mempool.AddPending(topic.Topic, topic.Duration); err != nil {
				topic.ResponseCh <- s.GetErrorSubmitTopic(err.Error())
				continue
			}

			resMsg := fmt.Sprintf("pending opening success { topic: %s }", topic.Topic)

			topic.ResponseCh <- s.GetSuccessSubmitTopic(resMsg)

		case vote := <-s.RequestVoteCh:
			id, tx := vote.Fragmentation()

			if err := s.mempool.CommitTransaction(id, tx); err != nil {
				vote.ResponseCh <- s.GetErrorSubmitVote(err.Error())
				continue
			}

			vote.ResponseCh <- s.GetSuccessSubmitVote(vote.Hash.String())

		case pended := <-s.pendedCh:
			go s.createNewBlock(pended)

		case <-s.ExitSignalCh:
			log.Println("exit signal detected")
			break labelServer
		}
	}
}

func (s *BlockChainServer) startgRPCListener() {
	tn, tp := s.getTopicListenerOpts()
	go s.startTopicListener(tn, tp, s.ExitSignalCh)

	vn, vp := s.getVoteListenerOpts()
	go s.startVoteListener(vn, vp, s.ExitSignalCh)
}

func (s *BlockChainServer) getTopicListenerOpts() (network string, port uint16) {
	return s.BlockChainServerOpts.TopicgRPCNetwork, s.BlockChainServerOpts.TopicgRPCNetworkPort
}

func (s *BlockChainServer) getVoteListenerOpts() (network string, port uint16) {
	return s.BlockChainServerOpts.VotegRPCNetwork, s.BlockChainServerOpts.VotegRPCNetworkPort
}

func (s *BlockChainServer) createNewBlock(p *mempool.Pended) {
	// prevHeight := s.blockChain.Height()
	// prevHeader, _ := s.blockChain.GetHeader(prevHeight)

	currentProtoBlock := block.NewProtoBlock(p.GetPendingID(), p.GetTxx())
	// currentBlock := block.NewBlockFromPrevHeader(prevHeader, currentProtoBlock)

	// log.Printf(util.BlockString("PROTOBLOCK: %s | { BlockHash: %s, TxLength: %d }"),
	// 	p.GetPendingID(), currentProtoBlock.MerkleRoot, currentProtoBlock.Len())

	select {
	case s.protoBlockCh <- currentProtoBlock:
		return
	default:
		log.Printf(util.BlockChainString("failed to push block %s: block channel is likely full or closed during send attempt"),
			currentProtoBlock.VotingID,
		)
		return
	}
}
