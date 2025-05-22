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
	RequestVoteCh chan *gRPC.PreTxVote
}

func NewBlockChainVoteServer(bufferSize int) *BlockChainVoteServer {
	return &BlockChainVoteServer{
		RequestVoteCh: make(chan *gRPC.PreTxVote, bufferSize),
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

	grpcServer := grpc.NewServer()

	vote_message.RegisterBlockchainVoteServiceServer(grpcServer, s)

	log.Printf(util.SystemString("SYSTEM: Vote gRPC listener opened { port: %d }"), port)

	if err := grpcServer.Serve(lis); err != nil {
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
