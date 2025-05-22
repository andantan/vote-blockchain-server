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

type BlockChainVoteListener struct {
	vote_message.UnimplementedBlockchainVoteServiceServer
	RequestVoteCh  chan *gRPC.PreTxVote
	ResponseVoteCh chan *gRPC.PostTxVote
}

func NewBlockChainVoteListener() *BlockChainVoteListener {
	return &BlockChainVoteListener{
		RequestVoteCh:  make(chan *gRPC.PreTxVote),
		ResponseVoteCh: make(chan *gRPC.PostTxVote),
	}
}

// gRPC
func (l *BlockChainVoteListener) SubmitVote(
	ctx context.Context, req *vote_message.VoteRequest,
) (*vote_message.VoteResponse, error) {
	preTxVote, err := gRPC.GetPreTxVote(req)

	if err != nil {
		return l.GetErrorSubmitVote(err.Error()).GetVoteResponse(), nil
	}

	l.RequestVoteCh <- preTxVote

	// Standby for reaching mempool: add Tx
	postTxVote := <-l.ResponseVoteCh

	return postTxVote.GetVoteResponse(), nil
}

func (l *BlockChainVoteListener) startVoteListener(network string, port uint16) {
	address := fmt.Sprintf(":%d", port) // ":port"

	lis, err := net.Listen(network, address)

	if err != nil {
		log.Fatalf("failed to listen on port 9000 (Vote): %v", err)
	}

	grpcServer := grpc.NewServer()

	vote_message.RegisterBlockchainVoteServiceServer(grpcServer, l)

	log.Printf("Vote gRPC listener opened (%d)", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server gRPC listener (Vote) over port 9000: %v", err)
	}
}

func (l *BlockChainVoteListener) GetSuccessSubmitVote(message string) *gRPC.PostTxVote {
	return gRPC.GetPostTxVote("SUCCESS", message, false)
}

func (l *BlockChainVoteListener) GetErrorSubmitVote(message string) *gRPC.PostTxVote {
	return gRPC.GetPostTxVote("ERROR", message, true)
}
