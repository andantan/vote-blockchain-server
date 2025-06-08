package main

import (
	"context"
	"log"
	"time"

	"github.com/andantan/vote-blockchain-server/client/gRPC_client/topic_client/topic_message"
	"github.com/andantan/vote-blockchain-server/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.NewClient(
		":9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	defer conn.Close()

	c := topic_message.NewBlockchainTopicServiceClient(conn)

	topics := []string{"2025 대선", "2025 보건의료 여론조사", "법률개정안 찬반 투표", "상법개정안 시범 기간 조사"}
	// topics := []string{"2025 대선", "2025 보건의료 여론조사", "법률개정안 찬반 투표"}
	// topics := []string{"2025 대선", "2025 보건의료 여론조사"}
	// topics := []string{"2025 대선"}
	for _, t := range topics {
		topic := topic_message.TopicRequest{
			Topic:    t,
			Duration: 5,
		}

		response, err := c.SubmitTopic(context.Background(), &topic)

		if err != nil {
			log.Fatalf("error when calling SubmitTopic: %s", err)
		}

		log.Printf("Response from server: %+v\n", response)
		r := util.RandRange(5, 20)
		time.Sleep(time.Duration(r) * time.Second)
	}
}
