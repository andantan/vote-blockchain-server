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
	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_submit_message"
	"github.com/andantan/vote-blockchain-server/util"
)

// gRPC Network and port options
//
// Network must be "tcp", "unix" or "unixpacket"
//
// Port, BufferSize must be between 0 and 65535
type VoteSubmitListenerOption struct {
	network           string
	port              uint16
	channelBufferSize uint16
}

func NewVoteSubmitListenerOption(network string, port, channelBufferSize uint16) *VoteSubmitListenerOption {
	return &VoteSubmitListenerOption{
		network:           network,
		port:              port,
		channelBufferSize: channelBufferSize,
	}
}

type VoteSubmitListener struct {
	vote_submit_message.UnimplementedBlockchainVoteSubmitServiceServer
	option       *VoteSubmitListenerOption
	ctx          context.Context
	cancel       context.CancelFunc
	grpcServer   *grpc.Server
	voteSubmitCh chan *gRPC.VoteSubmit
	exitCh       chan uint8
}

func NewVoteSubmitListener(exitCh chan uint8) *VoteSubmitListener {
	__cfg := config.GetGrpcVoteSubmitListenerConfiguration()
	__sys_channel_size := config.GetChannelBufferSizeSystemConfiguration()

	ctx, cancel := context.WithCancel(context.Background())

	opts := NewVoteSubmitListenerOption(
		__cfg.Network,
		__cfg.Port,
		__sys_channel_size.GrpcVoteProposalChannelBufferSize,
	)

	return &VoteSubmitListener{
		option:       opts,
		ctx:          ctx,
		cancel:       cancel,
		voteSubmitCh: make(chan *gRPC.VoteSubmit, opts.channelBufferSize),
		exitCh:       exitCh,
	}
}

func (li *VoteSubmitListener) SetGrpcServer(s *grpc.Server) {
	li.grpcServer = s
}

func (li *VoteSubmitListener) Consume() chan *gRPC.VoteSubmit {
	return li.voteSubmitCh
}

// gRPC
func (listener *VoteSubmitListener) SubmitBallotTransaction(
	ctx context.Context, req *vote_submit_message.SubmitBallotTransactionRequest,
) (*vote_submit_message.SubmitBallotTransactionResponse, error) {
	vs, err := gRPC.NewVoteSubmit(req)

	if err != nil {
		return gRPC.NewErrorVoteSubmitResponse(err).GetVoteResponse(), nil
	}

	ResponseCh := make(chan *gRPC.VoteSubmitResponse, 1)
	defer close(ResponseCh)

	vs.ResponseCh = ResponseCh

	listener.voteSubmitCh <- vs

	// Standby for reaching mempool: add Tx
	vsr := <-ResponseCh

	return vsr.GetVoteResponse(), nil

}

func (listener *VoteSubmitListener) Start(wg *sync.WaitGroup) {
	address := fmt.Sprintf(":%d", listener.option.port) // ":port"

	lis, err := net.Listen(listener.option.network, address)

	if err != nil {
		log.Printf(util.FatalString("failed to listen on port 9000 (Vote): %v"), err)
		listener.exitCh <- 0
	}

	listener.grpcServer = grpc.NewServer()

	go listener.stopVoteListener()

	vote_submit_message.RegisterBlockchainVoteSubmitServiceServer(listener.grpcServer, listener)

	log.Printf(util.SystemString("SYSTEM: Vote submit gRPC listener opened { port: %d }"), listener.option.port)
	wg.Done()

	if err := listener.grpcServer.Serve(lis); err != nil {
		log.Printf(util.FatalString("failed to server gRPC listener (VoteSubmit) over port %d: %v"), listener.option.port, err)
		listener.exitCh <- 0
	}
}

func (listener *VoteSubmitListener) stopVoteListener() {
	defer close(listener.voteSubmitCh)

	<-listener.ctx.Done()
	log.Println(util.SystemString("SYSTEM: Vote submit gRPC server received shutdown signal. Gracefully stopping..."))
	listener.grpcServer.GracefulStop()
	log.Println(util.SystemString("SYSTEM: Vote submit gRPC server stopped"))
}

func (listener *VoteSubmitListener) Shutdown() {
	log.Println(util.SystemString("SYSTEM: Requesting VoteSubmitListener to stop"))
	listener.cancel()
	log.Println(util.SystemString("SYSTEM: VoteSubmitListener has stopped"))
}
