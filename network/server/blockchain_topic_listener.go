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
	RequestTopicCh  chan *gRPC.PreTxTopic // TODO Request, Response channel
	ResponseTopicCh chan *gRPC.PostTxTopic
}

func NewBlockChainTopicListener() *BlockChainTopicListener {
	return &BlockChainTopicListener{
		RequestTopicCh:  make(chan *gRPC.PreTxTopic),
		ResponseTopicCh: make(chan *gRPC.PostTxTopic),
	}
}

// gRPC
func (l *BlockChainTopicListener) SubmitTopic(
	ctx context.Context, req *topic_message.TopicRequest,
) (*topic_message.TopicResponse, error) {
	l.RequestTopicCh <- gRPC.GetPreTxTopic(req)

	postTxTopic := <-l.ResponseTopicCh

	return postTxTopic.GetTopicResponse(), nil
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

func (l *BlockChainTopicListener) GetSuccessSubmitTopic(message string) *gRPC.PostTxTopic {
	return gRPC.GetPostTxTopic("SUCCESS", message, false)
}

func (l *BlockChainTopicListener) GetErrorSubmitTopic(message string) *gRPC.PostTxTopic {
	return gRPC.GetPostTxTopic("ERROR", message, true)
}

func (l *BlockChainTopicListener) GetFailedSubmitTopic(message string) *gRPC.PostTxTopic {
	return gRPC.GetPostTxTopic("FAILED", message, false)
}
