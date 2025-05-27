package blockchain

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	BLOCK_REQUEST_BUFFER_SIZE = 16
)

type BlockChain struct {
	mu      sync.RWMutex
	headers []*block.Header

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	blockCh chan *block.Block
}

func NewBlockChain() *BlockChain {
	ctx, cancel := context.WithCancel(context.Background())

	bc := &BlockChain{
		headers: []*block.Header{},
		ctx:     ctx,
		cancel:  cancel,
		wg:      &sync.WaitGroup{},
		blockCh: make(chan *block.Block, BLOCK_REQUEST_BUFFER_SIZE),
	}

	// sync go-routine Activate()
	bc.wg.Add(1)

	go bc.Activate()

	log.Printf(
		util.BlockChainString("BLOCKCHAIN: activate blockchain  { BLOCK_REQUEST_BUFFER_SIZE: %d }"),
		BLOCK_REQUEST_BUFFER_SIZE,
	)

	return bc
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

func (bc *BlockChain) Produce() chan<- *block.Block {
	return bc.blockCh
}

func (bc *BlockChain) Activate() {
	defer close(bc.blockCh)
	defer bc.wg.Done()

	for {
		select {
		case <-bc.ctx.Done():
			// TODO something...
			return
		case newBlock := <-bc.blockCh:
			log.Printf(util.BlockChainString("BLOCKCHAIN: received block %s | { BlockHash: %s, TxLength: %d }"),
				newBlock.VotingID, newBlock.BlockHash, len(newBlock.Transactions))
		}
	}
}

func (bc *BlockChain) Stop() {
	log.Println(util.SystemString("BlockChainServer: Stopping..."))

	bc.cancel()

	bc.wg.Wait()

	log.Println(util.SystemString("BlockChainServer: Stopped."))
}
