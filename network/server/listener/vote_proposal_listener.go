package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/andantan/vote-blockchain-server/network/gRPC"
	"github.com/andantan/vote-blockchain-server/network/gRPC/topic_message"
	"github.com/andantan/vote-blockchain-server/util"
)

// gRPC Network and port options
//
// Network must be "tcp", "unix" or "unixpacket"
//
// Port, BufferSize must be between 0 and 65535
type VoteProposalListenerOption struct {
	network           string
	port              uint16
	channelBufferSize uint16
}

func NewVoteProposalListenerOption(network string, port, channelBufferSize uint16) *VoteProposalListenerOption {
	return &VoteProposalListenerOption{
		network:           network,
		port:              port,
		channelBufferSize: channelBufferSize,
	}
}

type VoteProposalListener struct {
	*VoteProposalListenerOption
	topic_message.UnimplementedBlockchainTopicServiceServer
	ctx            context.Context
	cancel         context.CancelFunc
	grpcServer     *grpc.Server
	voteProposalCh chan *gRPC.VoteProposal
	exitCh         chan uint8
}

func NewVoteProposalListener(opts *VoteProposalListenerOption, exitCh chan uint8) *VoteProposalListener {
	ctx, cancel := context.WithCancel(context.Background())

	return &VoteProposalListener{
		VoteProposalListenerOption: opts,
		ctx:                        ctx,
		cancel:                     cancel,
		voteProposalCh:             make(chan *gRPC.VoteProposal, opts.channelBufferSize),
		exitCh:                     exitCh,
	}
}

func (li *VoteProposalListener) SetGrpcServer(s *grpc.Server) {
	li.grpcServer = s
}

func (li *VoteProposalListener) Consume() chan *gRPC.VoteProposal {
	return li.voteProposalCh
}

// gRPC
func (listener *VoteProposalListener) SubmitTopic(
	ctx context.Context, req *topic_message.TopicRequest,
) (*topic_message.TopicResponse, error) {
	ResponseCh := make(chan *gRPC.VoteProposalResponse, 1)
	defer close(ResponseCh)

	vp := gRPC.NewVoteProposal(req)
	vp.ResponseCh = ResponseCh

	listener.voteProposalCh <- vp

	// Standby for reach & validate to mempool: pending
	vpr := <-ResponseCh

	return vpr.GetTopicResponse(), nil
}

func (listener *VoteProposalListener) Start(wg *sync.WaitGroup) {
	address := fmt.Sprintf(":%d", listener.port) // ":port"

	lis, err := net.Listen(listener.network, address)

	if err != nil {
		log.Printf(util.FatalString("failed to listen on port 9001 (Topic): %v"), err)
		listener.exitCh <- 0
	}

	listener.grpcServer = grpc.NewServer()

	go listener.stopTopicListener()

	topic_message.RegisterBlockchainTopicServiceServer(listener.grpcServer, listener)

	log.Printf(util.SystemString("SYSTEM: Topic gRPC listener opened { port: %d }"), listener.port)
	wg.Done()

	if err := listener.grpcServer.Serve(lis); err != nil {
		log.Printf(util.FatalString("failed to server gRPC listener (Topic) over port %d: %v"), listener.port, err)
		listener.exitCh <- 0
	}
}

func (listener *VoteProposalListener) stopTopicListener() {
	defer close(listener.voteProposalCh)

	<-listener.ctx.Done()
	log.Println(util.SystemString("SYSTEM: Topic gRPC server received shutdown signal. Gracefully stopping..."))
	listener.grpcServer.GracefulStop()
	log.Println(util.SystemString("SYSTEM: Topic gRPC server stopped"))
}

func (listener *VoteProposalListener) Shutdown() {
	log.Println(util.SystemString("SYSTEM: Requesting BlockChainTopicServer to stop"))
	listener.cancel()
	log.Println(util.SystemString("SYSTEM: BlockChainTopicServer has stopped"))
}
