package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/andantan/vote-blockchain-server/network/gRPC"
	"github.com/andantan/vote-blockchain-server/network/gRPC/topic_message"
	"google.golang.org/grpc"
)

type BlockChainTopicServer struct {
	topic_message.UnimplementedBlockchainTopicServiceServer
	RequestTopicCh  chan *gRPC.PreTxTopic // TODO Request, Response channel
	ResponseTopicCh chan *gRPC.PostTxTopic
}

func NewBlockChainTopicServer() *BlockChainTopicServer {
	return &BlockChainTopicServer{
		RequestTopicCh:  make(chan *gRPC.PreTxTopic),
		ResponseTopicCh: make(chan *gRPC.PostTxTopic),
	}
}

// gRPC
func (s *BlockChainTopicServer) SubmitTopic(
	ctx context.Context, req *topic_message.TopicRequest,
) (*topic_message.TopicResponse, error) {
	s.RequestTopicCh <- gRPC.GetPreTxTopic(req)

	// Standby for reaching mempool: pending
	postTxTopic := <-s.ResponseTopicCh

	return postTxTopic.GetTopicResponse(), nil
}

func (s *BlockChainTopicServer) startTopicListener(network string, port uint16) {
	address := fmt.Sprintf(":%d", port) // ":port"

	lis, err := net.Listen(network, address)

	if err != nil {
		log.Fatalf("failed to listen on port 9001 (Topic): %v", err)
	}

	grpcServer := grpc.NewServer()

	topic_message.RegisterBlockchainTopicServiceServer(grpcServer, s)

	log.Printf("Topic gRPC listener opened (%d)", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server gRPC listener (Topic) over port 9001: %v", err)
	}
}

func (s *BlockChainTopicServer) GetSuccessSubmitTopic(message string) *gRPC.PostTxTopic {
	return gRPC.GetPostTxTopic("SUCCESS", message, false)
}

func (s *BlockChainTopicServer) GetErrorSubmitTopic(message string) *gRPC.PostTxTopic {
	return gRPC.GetPostTxTopic("ERROR", message, true)
}
