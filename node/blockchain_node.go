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
	TestBlockTime = 30 * time.Second
	TestMaxTxSize = uint32(600)
)

const (
	STORE_BASE_DIR   = "./"
	STORE_BLOCKS_DIR = "blocks"
)

const (
	VOTE_PROPOSAL_NETWORK             = "tcp"
	VOTE_PROPOSAL_PORT                = uint16(9000)
	VOTE_PROPOSAL_CHANNEL_BUFFER_SIZE = uint16(256)
)

const (
	VOTE_SUBMIT_NETWORK             = "tcp"
	VOTE_SUBMIT_PORT                = uint16(9001)
	VOTE_SUBMIT_CHANNEL_BUFFER_SIZE = uint16(2048)
)

func Start(wg *sync.WaitGroup) {
	defer wg.Done()

	listenerOption := server.NewListenerOption()
	listenerOption.SetVoteProposalListenerOption(VOTE_PROPOSAL_NETWORK, VOTE_PROPOSAL_PORT, VOTE_PROPOSAL_CHANNEL_BUFFER_SIZE)
	listenerOption.SetVoteSubmitListenerOption(VOTE_SUBMIT_NETWORK, VOTE_SUBMIT_PORT, VOTE_SUBMIT_CHANNEL_BUFFER_SIZE)

	blockOption := server.NewBlockOption(TestBlockTime, TestMaxTxSize)
	// blockOption := server.NewBlockOption(BlockTime, MaxTxSize)
	storeOption := server.NewStoreOption(STORE_BASE_DIR, STORE_BLOCKS_DIR)

	serverOption := server.NewBlockChainServerOpts()
	serverOption.SetListenerOption(listenerOption)
	serverOption.SetBlockOption(blockOption)
	serverOption.SetStoreOption(storeOption)

	blockChainServer := server.NewBlockChainServer(serverOption)

	blockChainServer.Start()
}
