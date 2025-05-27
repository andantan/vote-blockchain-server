package main

import (
	"sync"

	"github.com/andantan/vote-blockchain-server/node"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// quitch := make(chan struct{})

	go node.Start(wg)
	wg.Wait()
	//<-quitch
}
