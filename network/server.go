package network

import (
	"context"
	"log"
	"math/rand"

	"github.com/andantan/vote-blockchain-server/vote"
)

type Server struct {
	vote.UnimplementedBlockchainServiceServer
}

func (s *Server) SubmitVote(ctx context.Context, req *vote.VoteRequest) (*vote.VoteResponse, error) {
	log.Printf("Received vote from client: %+v\n", req)

	return &vote.VoteResponse{
		BlockHeight: int64(rand.Intn(10000)),
	}, nil
}
