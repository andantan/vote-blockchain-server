package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

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
	*VoteSubmitListenerOption
	vote_submit_message.UnimplementedBlockchainVoteSubmitServiceServer
	ctx          context.Context
	cancel       context.CancelFunc
	grpcServer   *grpc.Server
	voteSubmitCh chan *gRPC.VoteSubmit
	exitCh       chan uint8
}

func NewVoteSubmitListener(opts *VoteSubmitListenerOption, exitCh chan uint8) *VoteSubmitListener {
	ctx, cancel := context.WithCancel(context.Background())

	return &VoteSubmitListener{
		VoteSubmitListenerOption: opts,
		ctx:                      ctx,
		cancel:                   cancel,
		voteSubmitCh:             make(chan *gRPC.VoteSubmit, opts.channelBufferSize),
		exitCh:                   exitCh,
	}
}

func (li *VoteSubmitListener) SetGrpcServer(s *grpc.Server) {
	li.grpcServer = s
}

func (li *VoteSubmitListener) Consume() chan *gRPC.VoteSubmit {
	return li.voteSubmitCh
}

// gRPC
func (listener *VoteSubmitListener) SubmitVote(
	ctx context.Context, req *vote_submit_message.VoteSubmitRequest,
) (*vote_submit_message.VoteSubmitResponse, error) {
	vs, err := gRPC.NewVoteSubmit(req)

	if err != nil {
		return gRPC.GetErrorSubmitVote(err.Error()).GetVoteResponse(), nil
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
	address := fmt.Sprintf(":%d", listener.port) // ":port"

	lis, err := net.Listen(listener.network, address)

	if err != nil {
		log.Printf(util.FatalString("failed to listen on port 9000 (Vote): %v"), err)
		listener.exitCh <- 0
	}

	listener.grpcServer = grpc.NewServer()

	go listener.stopVoteListener()

	vote_submit_message.RegisterBlockchainVoteSubmitServiceServer(listener.grpcServer, listener)

	log.Printf(util.SystemString("SYSTEM: Vote gRPC listener opened { port: %d }"), listener.port)
	wg.Done()

	if err := listener.grpcServer.Serve(lis); err != nil {
		log.Printf(util.FatalString("failed to server gRPC listener (Vote) over port %d: %v"), listener.port, err)
		listener.exitCh <- 0
	}
}

func (listener *VoteSubmitListener) stopVoteListener() {
	defer close(listener.voteSubmitCh)

	<-listener.ctx.Done()
	log.Println(util.SystemString("SYSTEM: Vote gRPC server received shutdown signal. Gracefully stopping..."))
	listener.grpcServer.GracefulStop()
	log.Println(util.SystemString("SYSTEM: Vote gRPC server stopped"))
}

func (listener *VoteSubmitListener) Shutdown() {
	log.Println(util.SystemString("SYSTEM: Requesting BlockChainVoteServer to stop"))
	listener.cancel()
	log.Println(util.SystemString("SYSTEM: BlockChainVoteServer has stopped"))
}
