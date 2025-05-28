package node

import (
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/network/server"
)

const (
	BlockTime = 30 * time.Second
	MaxTxSize = uint32(300)
)

const (
	TestBlockTime = 10 * time.Second
	TestMaxTxSize = uint32(50)
)

func Start(wg *sync.WaitGroup) {
	defer wg.Done()

	BCopts := server.NewBlockChainServerOpts()
	BCopts.SetTopicOptions("tcp", uint16(9000))
	BCopts.SetVoteOptions("tcp", uint16(9001))
	BCopts.SetControllOptions(BlockTime, MaxTxSize)
	// BCopts.SetControllOptions(TestBlockTime, TestMaxTxSize)
	blockChainServer := server.NewBlockChainServer(BCopts)

	blockChainServer.Start()

}
