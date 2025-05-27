package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/andantan/vote-blockchain-server/network/gRPC"
	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_message"
	"github.com/andantan/vote-blockchain-server/util"
)

type BlockChainVoteServer struct {
	vote_message.UnimplementedBlockchainVoteServiceServer
	grpcServer    *grpc.Server
	RequestVoteCh chan *gRPC.PreTxVote
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewBlockChainVoteServer(bufferSize int) *BlockChainVoteServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &BlockChainVoteServer{
		RequestVoteCh: make(chan *gRPC.PreTxVote, bufferSize),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// gRPC
func (s *BlockChainVoteServer) SubmitVote(
	ctx context.Context, req *vote_message.VoteRequest,
) (*vote_message.VoteResponse, error) {
	preTxVote, err := gRPC.GetPreTxVote(req)

	if err != nil {
		return s.GetErrorSubmitVote(err.Error()).GetVoteResponse(), nil
	}

	ResponseCh := make(chan *gRPC.PostTxVote, 1)
	defer close(ResponseCh)

	preTxVote.ResponseCh = ResponseCh

	s.RequestVoteCh <- preTxVote

	// Standby for reaching mempool: add Tx
	postTxVote := <-ResponseCh

	return postTxVote.GetVoteResponse(), nil

}

func (s *BlockChainVoteServer) startVoteListener(network string, port uint16, exitCh chan<- uint8) {
	address := fmt.Sprintf(":%d", port) // ":port"

	lis, err := net.Listen(network, address)

	if err != nil {
		log.Printf(util.FatalString("failed to listen on port 9000 (Vote): %v"), err)
		exitCh <- EXIT_SIGNAL
	}

	s.grpcServer = grpc.NewServer()

	go s.stopVoteListener()

	vote_message.RegisterBlockchainVoteServiceServer(s.grpcServer, s)

	log.Printf(util.SystemString("SYSTEM: Vote gRPC listener opened { port: %d }"), port)

	if err := s.grpcServer.Serve(lis); err != nil {
		log.Printf(util.FatalString("failed to server gRPC listener (Vote) over port %d: %v"), port, err)
		exitCh <- EXIT_SIGNAL
	}
}

func (s *BlockChainVoteServer) GetSuccessSubmitVote(message string) *gRPC.PostTxVote {
	return gRPC.GetPostTxVote("SUCCESS", message, true)
}

func (s *BlockChainVoteServer) GetErrorSubmitVote(message string) *gRPC.PostTxVote {
	return gRPC.GetPostTxVote("ERROR", message, false)
}

func (s *BlockChainVoteServer) stopVoteListener() {
	defer close(s.RequestVoteCh)

	<-s.ctx.Done()
	log.Println(util.SystemString("SYSTEM: Vote gRPC server received shutdown signal. Gracefully stopping..."))
	s.grpcServer.GracefulStop()
	log.Println(util.SystemString("SYSTEM: Vote gRPC server stopped"))
}

func (s *BlockChainVoteServer) Shutdown() {
	log.Println(util.SystemString("SYSTEM: Requesting BlockChainVoteServer to stop"))
	s.cancel()
	log.Println(util.SystemString("SYSTEM: BlockChainVoteServer has stopped"))
}
