package node

import (
	"time"

	"github.com/andantan/vote-blockchain-server/network/server"
)

const (
	BlockTime = 5 * time.Second
	MaxTxSize = uint32(50000)
)

func Start() {
	BCopts := server.NewBlockChainServerOpts()
	BCopts.SetTopicOptions("tcp", uint16(9000))
	BCopts.SetVoteOptions("tcp", uint16(9001))
	BCopts.SetControllOptions(BlockTime, MaxTxSize)

	blockChainServer := server.NewBlockChainServer(BCopts)

	blockChainServer.Start()
}
