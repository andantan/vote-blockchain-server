package node

import (
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/network/server"
)

const (
	BlockTime = 1 * time.Minute
	MaxTxSize = uint32(320)
)

const (
	TestBlockTime = 5 * time.Second
	TestMaxTxSize = uint32(200)
)

const (
	STORE_BASE_DIR   = "./"
	STORE_BLOCKS_DIR = "blocks"
)

func Start(wg *sync.WaitGroup) {
	defer wg.Done()

	BCopts := server.NewBlockChainServerOpts()
	BCopts.SetTopicOptions("tcp", uint16(9000))
	BCopts.SetVoteOptions("tcp", uint16(9001))
	// BCopts.SetControllOptions(BlockTime, MaxTxSize)
	BCopts.SetControllOptions(TestBlockTime, TestMaxTxSize)
	BCopts.SetStorerDirectoryOptions(STORE_BASE_DIR, STORE_BLOCKS_DIR)
	blockChainServer := server.NewBlockChainServer(BCopts)

	blockChainServer.Start()

}
