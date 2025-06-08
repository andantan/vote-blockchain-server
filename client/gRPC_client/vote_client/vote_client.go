package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/andantan/vote-blockchain-server/client/gRPC_client/vote_client/vote_message"
	"github.com/andantan/vote-blockchain-server/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.NewClient(
		":9001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	defer conn.Close()

	c := vote_message.NewBlockchainVoteServiceClient(conn)

	for {
		topic := randTopic()

		vote := vote_message.VoteRequest{
			Hash:   util.RandomHash().String(),
			Option: randOpt(),
			Topic:  topic,
		}

		response, err := c.SubmitVote(context.Background(), &vote)

		if err != nil {
			log.Fatalf("error when calling SubmitVote: %s", err)
		}

		log.Printf("gRPC response success: %s | %t (%s | %s)", topic, response.Success, response.Status, response.Message)
		r := util.RandRange(10, 50)
		time.Sleep(time.Duration(r) * time.Millisecond)
		// time.Sleep(1 * time.Second)
	}
}

// var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var topics = []string{"2025 대선", "2025 보건의료 여론조사", "법률개정안 찬반 투표", "상법개정안 시범 기간 조사"}

// var topics = []string{"2025 대선", "2025 보건의료 여론조사", "법률개정안 찬반 투표"}

// var topics = []string{"2025 대선", "2025 보건의료 여론조사"}
var options = []rune("12345")

// func randSeq(n int) string {
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(b)
// }

func randTopic() string {
	return topics[rand.Intn(len(topics))]
}

func randOpt() string {
	return string(options[rand.Intn(len(options))])
}
