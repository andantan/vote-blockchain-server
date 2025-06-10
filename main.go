package main

import (
	"sync"

	"github.com/andantan/vote-blockchain-server/node"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go node.StartBlockChainNode(wg)
	go node.StartBlockChainExplorer(wg)
	wg.Wait()
}
