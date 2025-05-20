package main

import (
	"log"
	"time"

	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/network/server"
)

var defaultBlockTime = 5 * time.Second

func main() {
	blockChainServer := server.NewBlockChainServer()

	go blockChainServer.Start()
	ticker := time.NewTicker(defaultBlockTime)

	for {
		select {
		case vote := <-blockChainServer.VoteCh:
			log.Printf("received vote from client: %+v\n", vote)
		case <-ticker.C:
			block.CreateNewBlock()
		}

	}

}
