package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/andantan/vote-blockchain-server/network/gRPC"
	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_message"
	"google.golang.org/grpc"
)

type BlockChainVoteServer struct {
	vote_message.UnimplementedBlockchainVoteServiceServer
	RequestVoteCh  chan *gRPC.PreTxVote
	ResponseVoteCh chan *gRPC.PostTxVote
}

func NewBlockChainVoteServer() *BlockChainVoteServer {
	return &BlockChainVoteServer{
		RequestVoteCh:  make(chan *gRPC.PreTxVote),
		ResponseVoteCh: make(chan *gRPC.PostTxVote),
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

	// log.Printf("Submit vote %s|%s|%s\n", req.Hash, req.Option, req.Topic)
	s.RequestVoteCh <- preTxVote

	// Standby for reaching mempool: add Tx
	postTxVote := <-s.ResponseVoteCh
	return postTxVote.GetVoteResponse(), nil

}

func (s *BlockChainVoteServer) startVoteListener(network string, port uint16) {
	address := fmt.Sprintf(":%d", port) // ":port"

	lis, err := net.Listen(network, address)

	if err != nil {
		log.Fatalf("failed to listen on port 9000 (Vote): %v", err)
	}

	grpcServer := grpc.NewServer()

	vote_message.RegisterBlockchainVoteServiceServer(grpcServer, s)

	log.Printf("Vote gRPC listener opened (%d)", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server gRPC listener (Vote) over port 9000: %v", err)
	}
}

func (s *BlockChainVoteServer) GetSuccessSubmitVote(message string) *gRPC.PostTxVote {
	return gRPC.GetPostTxVote("SUCCESS", message, false)
}

func (s *BlockChainVoteServer) GetErrorSubmitVote(message string) *gRPC.PostTxVote {
	return gRPC.GetPostTxVote("ERROR", message, true)
}
