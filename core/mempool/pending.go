package mempool

import (
	"fmt"
	"log"
	"maps"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
)

type Pending struct {
	mu           sync.RWMutex
	transactions map[string]string
	txCache      map[string]uint8

	transactionCH        chan *transaction.Transaction
	pendingTime          time.Duration // Vote duration
	blockTime            time.Duration // Block Time (system)
	maxTransactionSize   uint32        // Tx size (system)
	scheduledBlockHeight []uint64      // Pended block heights
	pendingID            types.Topic   // TopicID
	closeCh              chan<- types.Topic
}

func NewPending(pendingTime, blockTime time.Duration,
	maxTxSize uint32, pendingId types.Topic, closeCh chan<- types.Topic) *Pending {
	return &Pending{
		transactions:         make(map[string]string),
		txCache:              make(map[string]uint8),
		transactionCH:        make(chan *transaction.Transaction),
		pendingTime:          pendingTime,
		blockTime:            blockTime,
		maxTransactionSize:   maxTxSize,
		scheduledBlockHeight: []uint64{},
		pendingID:            pendingId,
		closeCh:              closeCh,
	}
}

func (p *Pending) Len() int {
	p.mu.Lock()
	l := len(p.transactions)
	p.mu.Unlock()

	return l
}

func (p *Pending) Transactions() *map[string]string {
	m := make(map[string]string)

	maps.Copy(m, p.transactions)

	return &m
}

func (p *Pending) Activate() {
	blockTimer := time.NewTicker(p.blockTime)
	pendingTimer := time.NewTicker(p.pendingTime)

labelPending:
	for {
		select {
		case tx := <-p.transactionCH:
			p.commitTx(tx)

			txCount := p.Len()

			if p.maxTransactionSize <= uint32(txCount) {
				log.Printf(util.PendingString("PENDING: %s | Max transactions (%d) reached"),
					p.pendingID, txCount)

				p.flush()

				log.Printf(util.BlockString("Block: %s | Max tx size - New block created"), p.pendingID)
				log.Printf(util.PendingString("PENDING: %s | Transactions map cleared. New size: %d"),
					p.pendingID, p.Len())

				blockTimer.Reset(p.blockTime)
			}

		case <-blockTimer.C:
			p.flush()
			log.Printf(util.BlockString("Block: %s | Block timeout - new block created"), p.pendingID)

		case <-pendingTimer.C:
			p.flush()
			log.Printf(util.PendingString("Pending: %s | Pending is over"), p.pendingID)
			log.Printf(util.BlockString("Block: %s | Pending timeout -  New block created"), p.pendingID)
			log.Printf(util.PendingString("Pending: %s | close pending"), p.pendingID)

			p.closeCh <- p.pendingID

			break labelPending
		}
	}

	log.Printf(util.PendingString("Pending: %s | Activation exited."), p.pendingID)
}

func (p *Pending) PushTx(tx *transaction.Transaction) error {
	if p.collision(tx) {
		return fmt.Errorf("given tx (%s) is already commited", tx.GetHashString())
	}

	p.transactionCH <- tx

	return nil
}

func (p *Pending) collision(tx *transaction.Transaction) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.txCache[tx.GetHashString()]

	return ok
}

func (p *Pending) commitTx(tx *transaction.Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.commit(tx)
	p.cache(tx)

	log.Printf(util.CommitString("COMMIT: %s | New commit tx { hash: %s, option: %s }"),
		p.pendingID, tx.GetHashString(), tx.GetOption())
}

func (p *Pending) commit(tx *transaction.Transaction) {
	p.transactions[tx.GetHashString()] = tx.Serialize()
}

func (p *Pending) cache(tx *transaction.Transaction) {
	p.txCache[tx.GetHashString()] = 0
}

func (p *Pending) flush() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.transactions = make(map[string]string)
}
