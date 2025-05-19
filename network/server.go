package network

import (
	"context"
	"log"
	"math/rand"

	"github.com/andantan/vote-blockchain-server/protobuf"
)

type Server struct {
	protobuf.UnimplementedBlockchainServiceServer
}

func (s *Server) SubmitVote(ctx context.Context, req *protobuf.VoteRequest) (*protobuf.VoteResponse, error) {
	log.Printf("Received vote from client: %+v\n", req)

	return &protobuf.VoteResponse{
		BlockHeight: int64(rand.Intn(1000)),
	}, nil
}
