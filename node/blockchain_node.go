package node

import (
	"time"

	"github.com/andantan/vote-blockchain-server/network/server"
)

const (
	BlockTime = 10 * time.Second
	MaxTxSize = uint32(2000)
)

const (
	TestBlockTime = 10 * time.Second
	TestMaxTxSize = uint32(50)
)

func Start() {
	BCopts := server.NewBlockChainServerOpts()
	BCopts.SetTopicOptions("tcp", uint16(9000))
	BCopts.SetVoteOptions("tcp", uint16(9001))
	//BCopts.SetControllOptions(BlockTime, MaxTxSize)
	BCopts.SetControllOptions(TestBlockTime, TestMaxTxSize)
	blockChainServer := server.NewBlockChainServer(BCopts)

	blockChainServer.Start()
}
