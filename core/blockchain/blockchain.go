package blockchain

import (
	"fmt"
	"sync"

	"github.com/andantan/vote-blockchain-server/core/block"
)

type BlockChain struct {
	mu      sync.RWMutex
	headers []*block.Header
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		headers: []*block.Header{},
	}
}

func NewBlockChainWithGenesisBlock() *BlockChain {
	gb := block.GenesisBlock()
	bc := NewBlockChain()

	bc.attachBlock(gb)

	return bc
}

func (bc *BlockChain) GetHeader(height uint32) (*block.Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}
	bc.mu.Lock()
	defer bc.mu.Unlock()

	return bc.headers[int(height)], nil
}

// eg: headers [GenesisHeader, 1, 2, 3] => 4 len
// eg: headers [GenesisHeader, 1, 2, 3] => 3 height
func (bc *BlockChain) Height() uint32 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return uint32(len(bc.headers) - 1)
}

func (bc *BlockChain) attachBlock(b *block.Block) {
	bc.mu.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.mu.Unlock()
}
