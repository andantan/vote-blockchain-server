package node

import (
	"github.com/andantan/vote-blockchain-server/network/server"
)

func Start() {
	opts := server.NewServerOpts()
	opts.SetTopicOptions("tcp", uint16(9000))
	opts.SetVoteOptions("tcp", uint16(9001))

	blockChainServer := server.NewBlockChainServer(opts)

	blockChainServer.Start()
}
