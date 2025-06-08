package node

import (
	"sync"

	"github.com/andantan/vote-blockchain-server/network/server"
	SyncBlock "github.com/andantan/vote-blockchain-server/storage/sync"
)

func Start(wg *sync.WaitGroup) {
	defer wg.Done()

	validator := SyncBlock.NewValidator()
	validator.StartValidate()

	syncedHeaders := validator.GetSyncedBlockHeaders()
	blockChainServer := server.NewBlockChainServer(syncedHeaders)

	blockChainServer.Start()
}
