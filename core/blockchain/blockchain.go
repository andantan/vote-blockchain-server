package blockchain

import (
	"fmt"
	"log"
	"sync"

	"github.com/andantan/vote-blockchain-server/config"
	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/network/deliver"
	"github.com/andantan/vote-blockchain-server/storage/store"
	"github.com/andantan/vote-blockchain-server/util"
)

type BlockChain struct {
	mu           sync.RWMutex
	headers      []*block.Header
	wg           *sync.WaitGroup
	storer       *store.JsonStorer
	protoBlockCh chan *block.ProtoBlock
	eventDeliver *deliver.EventDeliver
}

func NewBlockChain(storer *store.JsonStorer, syncedHeader []*block.Header) *BlockChain {
	connectionUnicastBlockEventProtocol := config.GetEnvVar("CONNECTION_UNICAST_BLOCK_EVENT_PROTOCOL")
	connectionUnicastBlockEventAddress := config.GetEnvVar("CONNECTION_UNICAST_BLOCK_EVENT_ADDRESS")
	connectionUnicastBlockEventPort := config.GetIntEnvVar("CONNECTION_UNICAST_BLOCK_EVENT_PORT")
	connectionUnicastBlockEventEndpoint := config.GetEnvVar("CONNECTION_UNICAST_BLOCK_EVENT_ENDPOINT")

	bc := &BlockChain{
		wg:     &sync.WaitGroup{},
		storer: storer,
		eventDeliver: deliver.NewEventDeliver(
			connectionUnicastBlockEventProtocol,
			connectionUnicastBlockEventAddress,
			uint16(connectionUnicastBlockEventPort),
		),
	}

	bc.eventDeliver.SetCreatedBlockEventDeliver(connectionUnicastBlockEventEndpoint)

	log.Printf(util.DeliverString("DELIVER: Blockchain created block event deliver endpoint: %s"),
		bc.eventDeliver.CreatedBlockEventDeliver.GetUrl())

	if len(syncedHeader) != 0 {
		bc.setSyncedHeaders(syncedHeader)
	} else {
		bc.headers = []*block.Header{}
	}

	bc.setChannel()
	bc.wg.Add(1)

	go bc.Activate()

	return bc
}

func NewGenesisBlockChain(storer *store.JsonStorer) *BlockChain {
	gb := block.GenesisBlock()
	bc := NewBlockChain(storer, []*block.Header{})

	bc.attachBlock(gb)

	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block ID=%s"), gb.VotingID)
	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block MerkleRoot=%s"), gb.MerkleRoot.String())
	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block Height=%d"), gb.Height)
	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block PrevBlockHash=%s"), gb.PrevBlockHash.String())
	log.Printf(util.BlockChainString("BLOCKCHAIN: Genesis block BlockHash=%s"), gb.Hash().String())

	bc.storer.SaveBlock(gb)

	return bc
}

func (bc *BlockChain) setChannel() {
	log.Println(util.SystemString("SYSTEM: Blockchain setting channel..."))
	systemBlockPropaginateChannelBufferSize := config.GetIntEnvVar("SYSTEM_BLOCK_PROPAGINATE_CHANNEL_BUFFER_SIZE")

	bc.protoBlockCh = make(
		chan *block.ProtoBlock,
		systemBlockPropaginateChannelBufferSize,
	)

	log.Println(util.SystemString("SYSTEM: Blockchain block channel setting is done."))
}

func (bc *BlockChain) setSyncedHeaders(syncedHeader []*block.Header) {
	log.Printf(
		util.BlockChainString("SYNCHRONIZATION: Synchronize blockchain headers... | { Block-height: %d }"),
		len(syncedHeader)-1,
	)

	for _, header := range syncedHeader {
		bc.headers = append(bc.headers, header)

		log.Printf(util.BlockChainString("SYNCHRONIZATION: Header( 0x%s ) with Height( %d ) => "+util.YellowString("Synchronized")), header.Hash().String(), header.Height)
	}

	log.Println(util.BlockChainString("SYNCHRONIZATION: Blockchain headers synchronizion is done."))
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
| *H.VotingID      : %-80s
| *H.MerkleRoot    : %-80s
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

		bc.storer.SaveBlock(currentBlock)

		go bc.eventDeliver.UnicastCreatedBlockEvent(currentBlock)
	}

	log.Println(util.BlockChainString("BLOCKCHAIN: Block receiver and processor goroutine exited"))
}

func (bc *BlockChain) Shutdown() {
	log.Println(util.BlockChainString("BLOCKCHAIN: Initiating shutdown for BlockChain. Closing protoBlock channel"))
	close(bc.protoBlockCh)
	bc.wg.Wait()
	log.Println(util.BlockChainString("BLOCKCHAIN: BlockChain shutdown complete"))
}
