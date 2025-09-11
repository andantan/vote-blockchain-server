package node

import (
	"sync"

	"github.com/andantan/vote-blockchain-server/core/blockchain"
	"github.com/andantan/vote-blockchain-server/network/explorer"
)

func StartBlockChainExplorer(wg *sync.WaitGroup, chain *blockchain.BlockChain) {
	defer wg.Done()

	blockChainExplorer := explorer.NewBlockChainExplorer(chain)
	blockChainExplorer.Start()
}
