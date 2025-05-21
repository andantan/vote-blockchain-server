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

type BlockChainTopicListener struct {
	topic_message.UnimplementedBlockchainTopicServiceServer
	TopicCh chan gRPC.Topic
}

func (l *BlockChainTopicListener) SubmitTopic(
	ctx context.Context, req *topic_message.TopicRequest,
) (*topic_message.TopicResponse, error) {
	l.TopicCh <- gRPC.GetTopicFromTopicMessage(req)

	return &topic_message.TopicResponse{
		TopicId:       req.GetTopicId(),
		TopicDuration: req.GetTopicDuration(),
	}, nil
}

func (l *BlockChainTopicListener) startTopicListener(network string, port uint16) {
	address := fmt.Sprintf(":%d", port) // ":port"

	lis, err := net.Listen(network, address)

	if err != nil {
		log.Fatalf("failed to listen on port 9001 (Topic): %v", err)
	}

	grpcServer := grpc.NewServer()

	topic_message.RegisterBlockchainTopicServiceServer(grpcServer, l)

	log.Printf("Topic gRPC listener opened (%d)", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server gRPC listener (Topic) over port 9001: %v", err)
	}
}
