package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/andantan/vote-blockchain-server/config"
	"github.com/andantan/vote-blockchain-server/network/gRPC"
	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_proposal_message"
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
	vote_proposal_message.UnimplementedBlockchainVoteProposalServiceServer
	option         *VoteProposalListenerOption
	ctx            context.Context
	cancel         context.CancelFunc
	grpcServer     *grpc.Server
	voteProposalCh chan *gRPC.VoteProposal
	exitCh         chan uint8
}

func NewVoteProposalListener(exitCh chan uint8) *VoteProposalListener {
	__cfg := config.GetGrpcVoteProposalListenerConfiguration()
	__sys_channel_size := config.GetChannelBufferSizeSystemConfiguration()

	opts := NewVoteProposalListenerOption(
		__cfg.Network,
		__cfg.Port,
		__sys_channel_size.GrpcVoteProposalChannelBufferSize,
	)

	ctx, cancel := context.WithCancel(context.Background())

	return &VoteProposalListener{
		option:         opts,
		ctx:            ctx,
		cancel:         cancel,
		voteProposalCh: make(chan *gRPC.VoteProposal, opts.channelBufferSize),
		exitCh:         exitCh,
	}
}

func (li *VoteProposalListener) SetGrpcServer(s *grpc.Server) {
	li.grpcServer = s
}

func (li *VoteProposalListener) Consume() chan *gRPC.VoteProposal {
	return li.voteProposalCh
}

// gRPC
func (listener *VoteProposalListener) ProposalVote(
	ctx context.Context, req *vote_proposal_message.VoteProposalRequest,
) (*vote_proposal_message.VoteProposalResponse, error) {
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
	address := fmt.Sprintf(":%d", listener.option.port) // ":port"

	lis, err := net.Listen(listener.option.network, address)

	if err != nil {
		log.Printf(util.FatalString("failed to listen on port 9001 (Topic): %v"), err)
		listener.exitCh <- 0
	}

	listener.grpcServer = grpc.NewServer()

	go listener.stopTopicListener()

	vote_proposal_message.RegisterBlockchainVoteProposalServiceServer(listener.grpcServer, listener)

	log.Printf(util.SystemString("SYSTEM: Vote proposal gRPC listener opened { port: %d }"), listener.option.port)
	wg.Done()

	if err := listener.grpcServer.Serve(lis); err != nil {
		log.Printf(util.FatalString("failed to server gRPC listener (VoteProposal) over port %d: %v"), listener.option.port, err)
		listener.exitCh <- 0
	}
}

func (listener *VoteProposalListener) stopTopicListener() {
	defer close(listener.voteProposalCh)

	<-listener.ctx.Done()
	log.Println(util.SystemString("SYSTEM: Vote proposal gRPC server received shutdown signal. Gracefully stopping..."))
	listener.grpcServer.GracefulStop()
	log.Println(util.SystemString("SYSTEM: Vote proposal gRPC server stopped"))
}

func (listener *VoteProposalListener) Shutdown() {
	log.Println(util.SystemString("SYSTEM: Requesting VoteProposalListener to stop"))
	listener.cancel()
	log.Println(util.SystemString("SYSTEM: VoteProposalListener has stopped"))
}
