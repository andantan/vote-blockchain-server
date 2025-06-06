package node

import (
	"log"
	"sync"

	"github.com/andantan/vote-blockchain-server/network/server"
	SyncBlock "github.com/andantan/vote-blockchain-server/storage/sync"
	"github.com/andantan/vote-blockchain-server/util"
)

func Start(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println(util.SystemString("VALIDATE: Validate blockchain data"))

	validator := SyncBlock.NewValidator()
	validator.StartValidate()

	syncedHeaders := validator.GetSyncedBlockHeaders()

	blockChainServer := server.NewBlockChainServer(syncedHeaders)

	blockChainServer.Start()
}
