package main

import (
	"log"
	"sync"

	"github.com/andantan/vote-blockchain-server/config"
	"github.com/andantan/vote-blockchain-server/node"
)

func init() {
	log.Println("Initializing environment variables...")

	config.InitEnv()
}

func main() {
	wg := &sync.WaitGroup{}

	go node.StartBlockChainNode(wg)
	wg.Add(1)

	go node.StartBlockChainExplorer(wg)
	wg.Add(1)

	wg.Wait()
}
