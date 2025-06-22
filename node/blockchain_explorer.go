package node

import (
	"sync"

	"github.com/andantan/vote-blockchain-server/network/explorer"
)

func StartBlockChainExplorer(wg *sync.WaitGroup) {
	defer wg.Done()

	blockChainExplorer := explorer.NewBlockChainExplorer()
	blockChainExplorer.Start()
}
