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
	grpcServer     *grpc.Server
	RequestTopicCh chan *gRPC.PreTxTopic
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewBlockChainTopicServer(bufferSize int) *BlockChainTopicServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &BlockChainTopicServer{
		RequestTopicCh: make(chan *gRPC.PreTxTopic, bufferSize),
		ctx:            ctx,
		cancel:         cancel,
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

	s.grpcServer = grpc.NewServer()

	go s.stopTopicListener()

	topic_message.RegisterBlockchainTopicServiceServer(s.grpcServer, s)

	log.Printf(util.SystemString("SYSTEM: Topic gRPC listener opened { port: %d }"), port)

	if err := s.grpcServer.Serve(lis); err != nil {
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

func (s *BlockChainTopicServer) stopTopicListener() {
	defer close(s.RequestTopicCh)

	<-s.ctx.Done()
	log.Println(util.SystemString("SYSTEM: Topic gRPC server received shutdown signal. Gracefully stopping..."))
	s.grpcServer.GracefulStop()
	log.Println(util.SystemString("SYSTEM: Topic gRPC server stopped"))
}

func (s *BlockChainTopicServer) Shutdown() {
	log.Println(util.SystemString("SYSTEM: Requesting BlockChainTopicServer to stop"))
	s.cancel()
	log.Println(util.SystemString("SYSTEM: BlockChainTopicServer has stopped"))
}
