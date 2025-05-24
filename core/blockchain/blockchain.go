package blockchain

import (
	"sync"

	"github.com/andantan/vote-blockchain-server/core/block"
)

type BlockChain struct {
	mu     sync.RWMutex
	blocks []*block.Block
}
