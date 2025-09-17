package server

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/config"
	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/core/blockchain"
	"github.com/andantan/vote-blockchain-server/core/mempool"
	"github.com/andantan/vote-blockchain-server/network/deliver"
	"github.com/andantan/vote-blockchain-server/network/explorer"
	"github.com/andantan/vote-blockchain-server/network/gRPC"
	"github.com/andantan/vote-blockchain-server/network/server/listener"
	"github.com/andantan/vote-blockchain-server/storage/store"
	"github.com/andantan/vote-blockchain-server/util"
)

type BlockChainServer struct {
	*listener.VoteProposalListener
	*listener.VoteSubmitListener
	*listener.AdminCommandListener

	mempool    *mempool.MemPool
	blockChain *blockchain.BlockChain
	storer     *store.JsonStorer

	pendedCh     <-chan *mempool.Pended
	protoBlockCh chan<- *block.ProtoBlock
	exitSignalCh chan uint8

	eventDeliver *deliver.EventDeliver
}

func NewBlockChainServer(syncedHeaders []*block.Header) *BlockChainServer {
	server := &BlockChainServer{}
	server.exitSignalCh = make(chan uint8)

	server.initialize(syncedHeaders)

	log.Println(util.SystemString("SYSTEM: BlockChainServer engine generated"))
	log.Println(util.SystemString("SYSTEM: BlockChainServer initialization is done."))

	return server
}

func (s *BlockChainServer) initialize(syncedHeaders []*block.Header) {
	log.Println(util.SystemString("SYSTEM: BlockChainServer initialize..."))

	s.setgRPCServer()
	s.setMemPool()
	s.setStorer()
	s.setBlockChain(syncedHeaders)
	s.setChannel()
	s.setEventDeliver()

	log.Println(util.SystemString("SYSTEM: BlockChainServer initialization is done."))
}

func (s *BlockChainServer) setgRPCServer() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting gRPC server..."))

	s.VoteProposalListener = listener.NewVoteProposalListener(s.exitSignalCh)
	s.VoteSubmitListener = listener.NewVoteSubmitListener(s.exitSignalCh)
	s.AdminCommandListener = listener.NewAdminCommandListener(s.exitSignalCh)

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting gRPC server is done."))
}

func (s *BlockChainServer) setMemPool() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting memory pool..."))
	s.mempool = mempool.NewMemPool()
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting memory pool is done."))
}

func (s *BlockChainServer) setStorer() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer storer channel..."))

	s.storer = store.NewStore()

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting storer is done."))
}

func (s *BlockChainServer) setBlockChain(syncedHeaders []*block.Header) {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting blockchain..."))

	if len(syncedHeaders) == 0 {
		s.blockChain = blockchain.NewGenesisBlockChain(s.storer)

		if _, err := s.blockChain.GetHeader(0); err != nil {
			log.Fatalf(util.RedString("Genesis block initialization error: %s"), err.Error())
		}
	} else {
		s.blockChain = blockchain.NewBlockChain(s.storer, syncedHeaders)

		if _, err := s.blockChain.GetHeader(0); err != nil {
			log.Fatalf(util.RedString("Genesis block initialization error: %s"), err.Error())
		}
	}

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting blockchain is done."))
}

func (s *BlockChainServer) setChannel() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting channel..."))

	s.pendedCh = s.mempool.Consume()
	s.protoBlockCh = s.blockChain.Produce()

	log.Println(util.SystemString("SYSTEM: BlockChainServer setting channel is done."))
}

func (s *BlockChainServer) setEventDeliver() {
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting deliver..."))

	connectionUnicastPendingEventProtocol := config.GetEnvVar("CONNECTION_UNICAST_PENDING_EVENT_PROTOCOL")
	connectionUnicastPendingEventAddress := config.GetEnvVar("CONNECTION_UNICAST_PENDING_EVENT_ADDRESS")
	connectionUnicastPendingEventPort := config.GetIntEnvVar("CONNECTION_UNICAST_PENDING_EVENT_PORT")
	connectionUnicastPendingEventEndpoint := config.GetEnvVar("CONNECTION_UNICAST_PENDING_EVENT_ENDPOINT")

	s.eventDeliver = deliver.NewEventDeliver(
		connectionUnicastPendingEventProtocol,
		connectionUnicastPendingEventAddress,
		uint16(connectionUnicastPendingEventPort),
	)

	s.eventDeliver.SetExpirdPendingEventDeliver(connectionUnicastPendingEventEndpoint)

	log.Printf(util.DeliverString("DELIVER: BlockChainServer expired pending event deliver endpoint: %s"),
		s.eventDeliver.ExpiredPendingEventDeliver.GetUrl())
	log.Println(util.SystemString("SYSTEM: BlockChainServer setting deliver is done."))
}

func (s *BlockChainServer) Start() {
	s.startgRPCListener()

	explorer := explorer.NewBlockChainExplorer(s.blockChain, s.mempool)
	go explorer.Start()

	voteProposalCh := s.VoteProposalListener.Consume()
	voteSubmitCh := s.VoteSubmitListener.Consume()

labelServer:
	for {
		select {

		case proposal := <-voteProposalCh:
			// TODO: Make this tender
			if strings.Compare(string(proposal.Proposal), "exit") == 0 {
				log.Println(util.SystemString("============================== SHUTDOWN NODE =============================="))
				proposal.ResponseCh <- gRPC.NewSuccessVoteProposalResponse("shutdown", 1*time.Hour)

				s.VoteProposalListener.Shutdown()
				s.VoteSubmitListener.Shutdown()
				s.mempool.Shutdown()
				s.processClosedPendedCh()
				s.blockChain.Shutdown()
				s.storer.Shutdown()

				break labelServer
			}

			if err := s.mempool.AddPending(proposal.Proposal, proposal.Duration); err != nil {
				proposal.ResponseCh <- gRPC.NewErrorVoteProposalResponse(err)
				continue
			}

			proposal.ResponseCh <- gRPC.NewSuccessVoteProposalResponse(proposal.Proposal, proposal.Duration)

		case submit := <-voteSubmitCh:
			id, tx := submit.Fragmentation()

			if err := s.mempool.CommitTransaction(id, tx); err != nil {
				submit.ResponseCh <- gRPC.NewErrorVoteSubmitResponse(err)
				continue
			}

			submit.ResponseCh <- gRPC.NewSuccessVoteSubmitResponse(tx.Hash)

		case pended := <-s.pendedCh:
			if pended.IsExpired() {
				go s.eventDeliver.UnicastExpiredPendingEvent(pended)

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
	wg.Add(3)

	go s.VoteProposalListener.Start(wg)
	go s.VoteSubmitListener.Start(wg)
	go s.AdminCommandListener.Start(wg)

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

func (s *BlockChainServer) processClosedPendedCh() {
	log.Println(util.MemPoolString("MEMPOOL: Process closed tx channel"))

	for pended := range s.pendedCh {
		if pended.IsExpired() {
			log.Printf(" %+v | %+v\n", pended, pended.GetCachedOptions())

			continue
		}

		go s.createNewBlock(pended)
	}
}
