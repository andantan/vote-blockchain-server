package main

import "github.com/andantan/vote-blockchain-server/node"

func main() {
	quitch := make(chan int)

	go node.Start()

	<-quitch
}
