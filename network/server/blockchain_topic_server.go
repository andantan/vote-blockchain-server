package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/andantan/vote-blockchain-server/network/gRPC"
	"github.com/andantan/vote-blockchain-server/network/gRPC/topic_message"
	"github.com/andantan/vote-blockchain-server/util"
)

type BlockChainTopicServer struct {
	topic_message.UnimplementedBlockchainTopicServiceServer
	RequestTopicCh chan *gRPC.PreTxTopic
}

func NewBlockChainTopicServer(bufferSize int) *BlockChainTopicServer {
	return &BlockChainTopicServer{
		RequestTopicCh: make(chan *gRPC.PreTxTopic, bufferSize),
	}
}

// gRPC
func (s *BlockChainTopicServer) SubmitTopic(
	ctx context.Context, req *topic_message.TopicRequest,
) (*topic_message.TopicResponse, error) {
	ResponseCh := make(chan *gRPC.PostTxTopic, 1)
	defer close(ResponseCh)

	preTxTopic := gRPC.GetPreTxTopic(req)
	preTxTopic.ResponseCh = ResponseCh

	s.RequestTopicCh <- preTxTopic

	// Standby for reaching mempool: pending
	postTxTopic := <-ResponseCh

	return postTxTopic.GetTopicResponse(), nil
}

func (s *BlockChainTopicServer) startTopicListener(network string, port uint16, exitCh chan<- uint8) {
	address := fmt.Sprintf(":%d", port) // ":port"

	lis, err := net.Listen(network, address)

	if err != nil {
		log.Printf(util.FatalString("failed to listen on port 9001 (Topic): %v"), err)
		exitCh <- EXIT_SIGNAL
	}

	grpcServer := grpc.NewServer()

	topic_message.RegisterBlockchainTopicServiceServer(grpcServer, s)

	log.Printf(util.SystemString("SYSTEM: Topic gRPC listener opened { port: %d }"), port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf(util.FatalString("failed to server gRPC listener (Topic) over port %d: %v"), port, err)
		exitCh <- EXIT_SIGNAL
	}
}

func (s *BlockChainTopicServer) GetSuccessSubmitTopic(message string) *gRPC.PostTxTopic {
	return gRPC.GetPostTxTopic("SUCCESS", message, true)
}

func (s *BlockChainTopicServer) GetErrorSubmitTopic(message string) *gRPC.PostTxTopic {
	return gRPC.GetPostTxTopic("ERROR", message, false)
}
