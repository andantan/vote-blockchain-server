package server

import (
	"log"
	"time"

	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/network/gRPC"
)

const (
	BlockTime = 5 * time.Second
	MaxTxSize = uint32(50000)
)

// gRPC Network and port options
//
// Network must be "tcp", "unix" or "unixpacket"
//
// Port must be between 0 and 65535
type ServerOpts struct {
	TopicgRPCNetwork     string
	TopicgRPCNetworkPort uint16
	VotegRPCNetwork      string
	VotegRPCNetworkPort  uint16
}

func NewServerOpts() ServerOpts {
	return ServerOpts{}
}

func (o *ServerOpts) SetTopicOptions(topicNetwork string, topicNetworkPort uint16) {
	o.TopicgRPCNetwork = topicNetwork
	o.TopicgRPCNetworkPort = topicNetworkPort
}

func (o *ServerOpts) SetVoteOptions(voteNetwork string, voteNetworkPort uint16) {
	o.VotegRPCNetwork = voteNetwork
	o.VotegRPCNetworkPort = voteNetworkPort
}

type BlockChainServer struct {
	// vote_message.UnimplementedBlockchainVoteServiceServer
	// VoteCh chan gRPC.Vote
	ServerOpts
	*BlockChainVoteListener
	*BlockChainTopicListener
	ExitSignalCh chan uint8
}

func NewBlockChainServer(opts ServerOpts) *BlockChainServer {

	return &BlockChainServer{
		ServerOpts: opts,
		BlockChainVoteListener: &BlockChainVoteListener{
			VoteCh: make(chan gRPC.Vote),
		},
		BlockChainTopicListener: &BlockChainTopicListener{
			TopicCh: make(chan gRPC.Topic),
		},
		ExitSignalCh: make(chan uint8),
	}
}

func (s *BlockChainServer) Start() {
	ticker := time.NewTicker(BlockTime)

	go s.startTopicListener(s.getTopicOpts())
	go s.startVoteListener(s.getVoteOpts())

labelServer:
	for {
		select {
		case topic := <-s.TopicCh:
			log.Printf("received topic from client: %+v\n", topic)
		case vote := <-s.VoteCh:
			log.Printf("received vote from client: %+v\n", vote)
		case <-ticker.C:
			block.CreateNewBlock()
		case <-s.ExitSignalCh:
			log.Println("Exit signal detected")
			break labelServer
		}
	}
}

func (s *BlockChainServer) getTopicOpts() (network string, port uint16) {
	return s.ServerOpts.TopicgRPCNetwork, s.ServerOpts.TopicgRPCNetworkPort
}

func (s *BlockChainServer) getVoteOpts() (network string, port uint16) {
	return s.ServerOpts.VotegRPCNetwork, s.ServerOpts.VotegRPCNetworkPort
}
