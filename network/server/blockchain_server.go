package server

import (
	"log"
	"time"

	"github.com/andantan/vote-blockchain-server/core/mempool"
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
	mempool      *mempool.MemPool
	ExitSignalCh chan uint8
}

func NewBlockChainServer(opts BlockChainServerOpts) *BlockChainServer {

	return &BlockChainServer{
		BlockChainServerOpts:  opts,
		BlockChainVoteServer:  NewBlockChainVoteServer(),
		BlockChainTopicServer: NewBlockChainTopicServer(),
		mempool:               mempool.NewMemPool(opts.BlockTime, opts.MaxTxSize),
		ExitSignalCh:          make(chan uint8),
	}
}

func (s *BlockChainServer) Start() {

	go s.startTopicListener(s.getTopicListenerOpts())
	go s.startVoteListener(s.getVoteListenerOpts())

labelServer:
	for {
		select {
		case topic := <-s.RequestTopicCh:
			if err := s.mempool.AddPending(topic.Topic, topic.Duration); err != nil {
				s.ResponseTopicCh <- s.GetErrorSubmitTopic(err.Error())
				continue
			}

			s.ResponseTopicCh <- s.GetSuccessSubmitTopic("pending success (Topic)" + string(topic.Topic))

		case vote := <-s.RequestVoteCh:
			id, tx := vote.Fragmentation()

			if err := s.mempool.CommitTransaction(id, tx); err != nil {
				s.ResponseVoteCh <- s.GetErrorSubmitVote(err.Error())
				continue
			}

			s.ResponseVoteCh <- s.GetSuccessSubmitVote(vote.Hash.String())

		case <-s.ExitSignalCh:
			log.Println("exit signal detected")
			break labelServer
		}
	}
}

func (s *BlockChainServer) getTopicListenerOpts() (network string, port uint16) {
	return s.BlockChainServerOpts.TopicgRPCNetwork, s.BlockChainServerOpts.TopicgRPCNetworkPort
}

func (s *BlockChainServer) getVoteListenerOpts() (network string, port uint16) {
	return s.BlockChainServerOpts.VotegRPCNetwork, s.BlockChainServerOpts.VotegRPCNetworkPort
}
