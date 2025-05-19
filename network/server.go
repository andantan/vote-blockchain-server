package network

import (
	"context"
	"fmt"
	"math/rand"
	"net"

	"github.com/andantan/vote-blockchain-server/vote"
	"google.golang.org/grpc"
)

type Server struct {
	vote.UnimplementedBlockchainServiceServer
	VoteCh chan Vote
}

func NewServer() *Server {
	return &Server{
		VoteCh: make(chan Vote),
	}
}

// gRPC
func (s *Server) SubmitVote(ctx context.Context, req *vote.VoteRequest) (*vote.VoteResponse, error) {
	s.VoteCh <- GetVoteFromgRPCRequest(req)

	return &vote.VoteResponse{
		BlockHeight: int64(rand.Intn(10000)),
	}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":9000")

	if err != nil {
		return fmt.Errorf("failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()

	vote.RegisterBlockchainServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to server gRPC server over port 9000: %v", err)
	}

	return nil
}
