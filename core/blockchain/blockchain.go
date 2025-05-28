package blockchain

import (
	"fmt"
	"log"
	"sync"

	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	BLOCK_REQUEST_BUFFER_SIZE = 64
)

type BlockChain struct {
	mu           sync.RWMutex
	headers      []*block.Header
	wg           *sync.WaitGroup
	protoBlockCh chan *block.ProtoBlock
}

func NewBlockChain() *BlockChain {

	bc := &BlockChain{
		headers: []*block.Header{},
		wg:      &sync.WaitGroup{},
	}

	bc.setChannel()

	bc.wg.Add(1)

	go bc.Activate()

	return bc
}

func NewBlockChainWithGenesisBlock() *BlockChain {
	gb := block.GenesisBlock()
	bc := NewBlockChain()

	bc.attachBlock(gb)

	return bc
}

func (bc *BlockChain) setChannel() {
	log.Printf(
		util.SystemString("SYSTEM: Blockchain setting channel... | { BLOCK_REQUEST_BUFFER_SIZE: %d }"),
		BLOCK_REQUEST_BUFFER_SIZE,
	)

	bc.protoBlockCh = make(chan *block.ProtoBlock, BLOCK_REQUEST_BUFFER_SIZE)

	log.Println(util.SystemString("SYSTEM: Blockchain block channel setting is done."))
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

func (bc *BlockChain) Produce() chan<- *block.ProtoBlock {
	return bc.protoBlockCh
}

const (
	BLOCKCHAIN_NEW_CHAINED_BLOCK_LOG_MESSAGE = `BLOCKCHAIN: New block added to the chain.
--------------------------------------------------------------------------------------
| *H.Voting ID     : %-80s
| *H.Merkle Root   : %-80s
| *H.Height        : %-80d
| *H.PrevBlockHash : %-80s
| B.BlockHash      : %-80s
| B.TxLength       : %-80d
--------------------------------------------------------------------------------------
`
)

func (bc *BlockChain) Activate() {
	defer bc.wg.Done()

	log.Println(util.BlockChainString("BLOCKCHAIN: Starting block receiver and processor goroutine"))

	for pb := range bc.protoBlockCh {
		// log.Printf(util.BlockChainString("BLOCKCHAIN: Received protoBlock %s | { MerkleRoot: %s, TxxLength: %d }"),
		// 	pb.VotingID, pb.MerkleRoot, pb.Len())

		height := bc.Height()
		prevHeader, err := bc.GetHeader(height)

		if err != nil {
			log.Printf(util.BlockChainString("BLOCKCHAIN: GetHeader error (%s)"), err.Error())
			continue
		}

		currentBlock := block.NewBlockFromPrevHeader(prevHeader, pb)
		bc.attachBlock(currentBlock)

		log.Printf(util.BlockChainString(BLOCKCHAIN_NEW_CHAINED_BLOCK_LOG_MESSAGE),
			currentBlock.VotingID,
			currentBlock.MerkleRoot.String(),
			currentBlock.Height,
			currentBlock.PrevBlockHash.String(),
			currentBlock.BlockHash.String(),
			len(currentBlock.Transactions),
		)
	}

	log.Println(util.BlockChainString("BLOCKCHAIN: Block receiver and processor goroutine exited"))
}

func (bc *BlockChain) Shutdown() {
	log.Println(util.BlockChainString("BLOCKCHAIN: Initiating shutdown for BlockChain. Closing protoBlock channel"))
	close(bc.protoBlockCh)
	bc.wg.Wait()
	log.Println(util.BlockChainString("BLOCKCHAIN: BlockChain shutdown complete"))
}
