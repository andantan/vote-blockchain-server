package main

import (
	"sync"

	"github.com/andantan/vote-blockchain-server/node"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go node.Start(wg)
	wg.Wait()
}
