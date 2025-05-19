package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/andantan/vote-blockchain-server/vote"
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

	c := vote.NewBlockchainServiceClient(conn)

	for {
		vote := vote.VoteRequest{
			VoteHash:   randSeq(10),
			VoteOption: randSeq(10),
			ElectionId: randSeq(10),
		}

		response, err := c.SubmitVote(context.Background(), &vote)

		if err != nil {
			log.Fatalf("error when calling SubmitVote: %s", err)
		}

		log.Printf("Response from server: %+v\n", response)

		time.Sleep(50 * time.Millisecond)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
