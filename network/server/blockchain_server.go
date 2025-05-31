package server

import (
	"fmt"
	"log"
	"sync"

	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/core/blockchain"
	"github.com/andantan/vote-blockchain-server/core/mempool"
	"github.com/andantan/vote-blockchain-server/network/gRPC"
	"github.com/andantan/vote-blockchain-server/network/server/listener"
	"github.com/andantan/vote-blockchain-server/storage/store"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	EXIT_SIGNAL = iota
	CONVERT_SIGNAL
	STANDBY_SIGNAL
)

type BlockChainServer struct {
	serverOption BlockChainServerOpts

	*listener.VoteProposalListener
	*listener.VoteSubmitListener

	mempool    *mempool.MemPool
	blockChain *blockchain.BlockChain
	storer     *store.JsonStorer

	pendedCh     <-chan *mempool.Pended
	protoBlockCh chan<- *block.ProtoBlock
	exitSignalCh chan uint8
}

func NewBlockChainServer(options BlockChainServerOpts) *BlockChainServer {
	server := &BlockChainServer{serverOption: options}
	server.exitSignalCh = make(chan uint8)

	log.Println(util.SystemString("SYSTEM: BlockChainServer initialize..."))

	server.Initialize()

	log.Println(util.SystemString("SYSTEM: BlockChainServer engine generated"))
	log.Println(util.SystemString("SYSTEM: BlockChainServer initialization is done."))

	return server
}

func (s *BlockChainServer) Initialize() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer initialize..."))

	s.setgRPCServer()
	s.setMemPool()
	s.setStorer()
	s.setBlockChain()
	s.setChannel()

	log.Println(util.SystemString("SYSTEM: BlockChainServer initialization is done."))
}

func (s *BlockChainServer) setgRPCServer() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting gRPC server..."))

	s.VoteProposalListener = listener.NewVoteProposalListener(s.serverOption.VoteProposalListenerOption, s.exitSignalCh)
	s.VoteSubmitListener = listener.NewVoteSubmitListener(s.serverOption.VoteSubmitListenerOption, s.exitSignalCh)

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting gRPC server is done."))
}

func (s *BlockChainServer) setMemPool() {
	log.Printf(
		util.SystemString("SYSTEM: BlockChainServer setting memory pool... | { BlockTime: %s, MaxTxSize: %d }"),
		s.serverOption.blockTime,
		s.serverOption.maxTxSize,
	)

	s.mempool = mempool.NewMemPool(
		s.serverOption.blockTime,
		s.serverOption.maxTxSize,
	)

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting memory pool is done."))
}

func (s *BlockChainServer) setStorer() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer storer channel..."))

	s.storer = store.NewStore(
		s.serverOption.StoreBaseDirectory,
		s.serverOption.StoreBlocksDirectory,
	)

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting storer is done."))
}

func (s *BlockChainServer) setBlockChain() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting blockchain..."))

	s.blockChain = blockchain.NewBlockChainWithGenesisBlock(s.storer)

	if _, err := s.blockChain.GetHeader(0); err != nil {
		log.Fatalf(util.RedString("Genesis block initialization error: %s"), err.Error())
	}

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting blockchain is done."))
}

func (s *BlockChainServer) setChannel() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting channel..."))

	s.pendedCh = s.mempool.Consume()
	s.protoBlockCh = s.blockChain.Produce()

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting channel is done."))
}

func (s *BlockChainServer) Start() {
	s.startgRPCListener()

	voteProposalCh := s.VoteProposalListener.Consume()
	voteSubmitCh := s.VoteSubmitListener.Consume()

labelServer:
	for {
		select {
		case proposal := <-voteProposalCh:
			if err := s.mempool.AddPending(proposal.Topic, proposal.Duration); err != nil {
				proposal.ResponseCh <- gRPC.GetErrorVoteProposal(err.Error())
				continue
			}

			resMsg := fmt.Sprintf("pending opening success { topic: %s }", proposal.Topic)

			proposal.ResponseCh <- gRPC.GetSuccessVoteProposal(resMsg)

		case submit := <-voteSubmitCh:
			id, tx := submit.Fragmentation()

			if err := s.mempool.CommitTransaction(id, tx); err != nil {
				submit.ResponseCh <- gRPC.GetErrorSubmitVote(err.Error())
				continue
			}

			submit.ResponseCh <- gRPC.GetSuccessSubmitVote(submit.Hash.String())

		case pended := <-s.pendedCh:
			if pended.IsExpired() {
				log.Printf(" %+v | %+v\n", pended, pended.GetCachedOptions())

				continue
			}

			go s.createNewBlock(pended)

		case <-s.exitSignalCh:
			log.Println("exit signal detected")
			break labelServer
		}
	}
}

func (s *BlockChainServer) startgRPCListener() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go s.VoteProposalListener.Start(wg)
	go s.VoteSubmitListener.Start(wg)

	wg.Wait()
}

func (s *BlockChainServer) createNewBlock(p *mempool.Pended) {
	currentProtoBlock := block.NewProtoBlock(p.GetPendingID(), p.GetTxx())

	select {
	case s.protoBlockCh <- currentProtoBlock:
		return
	default:
		log.Printf(util.BlockChainString("failed to push block %s: block channel is likely full or closed during send attempt"),
			currentProtoBlock.VotingID)
		return
	}
}
