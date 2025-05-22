package main

import (
	"context"
	"log"
	"time"

	"github.com/andantan/vote-blockchain-server/client/topic_client/topic_message"
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

	for {
		topic := topic_message.TopicRequest{
			Topic:    "2025-경기도-철도공사",
			Duration: 7200,
		}

		response, err := c.SubmitTopic(context.Background(), &topic)

		if err != nil {
			log.Fatalf("error when calling SubmitTopic: %s", err)
		}

		log.Printf("Response from server: %+v\n", response)

		// time.Sleep(300 * time.Millisecond)
		time.Sleep(3 * time.Second)
	}
}
