package main

import (
	"flag"

	"github.com/andantan/vote-blockchain-server/impulse-client/client"
	"github.com/andantan/vote-blockchain-server/impulse-client/util"
)

func main() {
	var max int
	var registerMode bool

	flag.IntVar(&max, "max", 10, "Controls the number of operations. Its meaning changes based on '--register'.")
	flag.BoolVar(&registerMode, "register", false, "When true, enables user registration mode. Affects the meaning of '--max'.")
	flag.Parse()

	if registerMode {
		util.MakeRandomUsers(max)
		client.RegisterUsers()
	} else {
		openedVotes := client.BurstProposalClient(max)
		client.BurstSubmitClient(max, openedVotes)
	}

}
