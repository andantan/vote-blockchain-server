package server

import (
	"context"
	"fmt"
	"math/rand"
	"net"

	"github.com/andantan/vote-blockchain-server/network/gRPC"
	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_message"
	"google.golang.org/grpc"
)

type BlockChainServer struct {
	vote_message.UnimplementedBlockchainServiceServer
	VoteCh chan gRPC.Vote
}

func NewBlockChainServer() *BlockChainServer {
	return &BlockChainServer{
		VoteCh: make(chan gRPC.Vote),
	}
}

// gRPC
func (s *BlockChainServer) SubmitVote(ctx context.Context, req *vote_message.VoteRequest) (*vote_message.VoteResponse, error) {
	s.VoteCh <- gRPC.GetVoteFromVoteMessage(req)

	return &vote_message.VoteResponse{
		BlockHeight: int64(rand.Intn(10000)),
	}, nil
}

func (s *BlockChainServer) Start() error {
	lis, err := net.Listen("tcp", ":9000")

	if err != nil {
		return fmt.Errorf("failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()

	vote_message.RegisterBlockchainServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to server gRPC server over port 9000: %v", err)
	}

	return nil
}
