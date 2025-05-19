package main

import (
	"log"

	"github.com/andantan/vote-blockchain-server/network"
)

func main() {
	server := network.NewServer()

	go server.Start()

	for {
		grpc := <-server.VoteCh
		log.Printf("received vote from client: %+v\n", grpc)
	}

}
