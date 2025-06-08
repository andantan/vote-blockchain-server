package main

import (
	"flag"

	"github.com/andantan/vote-blockchain-server/impulse-client/client"
)

func main() {
	var max int
	flag.IntVar(&max, "max", 10, "Limit proposal entity")
	flag.Parse()

	client.BurstProposalClient(max)
	client.BurstSubmitClient(max)
}
