package server

import (
	"fmt"
	"log"
	"time"

	"github.com/andantan/vote-blockchain-server/core/mempool"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	EXIT_SIGNAL = iota
	CONVERT_SIGNAL
	STANDBY_SIGNAL
)

const (
	GRPC_REQUEST_BUFFER_SIZE = 128
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
	mempool  *mempool.MemPool
	pendedCh chan *mempool.Pended
	// chain        *blockchain.BlockChain
	ExitSignalCh chan uint8
}

func NewBlockChainServer(opts BlockChainServerOpts) *BlockChainServer {

	// Need genesisBlock
	return &BlockChainServer{
		BlockChainServerOpts:  opts,
		BlockChainVoteServer:  NewBlockChainVoteServer(GRPC_REQUEST_BUFFER_SIZE),
		BlockChainTopicServer: NewBlockChainTopicServer(GRPC_REQUEST_BUFFER_SIZE),
		mempool:               mempool.NewMemPool(opts.BlockTime, opts.MaxTxSize),
		pendedCh:              make(chan *mempool.Pended),
		ExitSignalCh:          make(chan uint8),
	}
}

func (s *BlockChainServer) Start() {
	s.startgRPCListener()
	s.mempool.SetChannel(s.pendedCh)

labelServer:
	for {
		select {
		case topic := <-s.RequestTopicCh:

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
			for _, v := range pended.GetTxx() {
				log.Printf(util.PendingString("PENDED: { %s }"), v.Serialize())
			}

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

// func (s *BlockChainServer) createNewBlock() {

// }
