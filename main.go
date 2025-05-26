package main

import (
	"log"

	"github.com/andantan/vote-blockchain-server/node"
	"github.com/google/gops/agent"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatalf("gops agent failed to start: %v", err)
	}

	quitch := make(chan int)

	go node.Start()

	<-quitch
}
