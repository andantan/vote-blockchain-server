package mempool

import (
	"fmt"
	"log"
	"sync"
	"time"

	"maps"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type Pending struct {
	mutex        sync.RWMutex
	transactions map[string]string

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
	p.mutex.Lock()
	l := len(p.transactions)
	p.mutex.Unlock()

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
				log.Printf("PENDING(%s) - Max transactions reached (%d). Creating new block by count.\n",
					p.pendingID, txCount)

				p.flush()

				log.Printf("New (%s) block created\n", p.pendingID)
				log.Printf("PENDING(%s) - Transactions map cleared. New size: %d\n", p.pendingID, p.Len())

				blockTimer.Reset(p.blockTime)
			}

		case <-blockTimer.C:
			p.flush()
			log.Printf("New (%s) block created\n", p.pendingID)

		case <-pendingTimer.C:
			p.flush()
			log.Printf("Pending(%s) over\n", p.pendingID)
			log.Printf("New block created(%s) by close pending\n", p.pendingID)

			p.closeCh <- p.pendingID

			break labelPending
		}
	}

	log.Printf("Pending(%s) activation exited.\n", p.pendingID)
}

func (p *Pending) PushTx(tx *transaction.Transaction) error {
	if p.collision(tx) {
		return fmt.Errorf("given tx (%s) is already commited", tx.GetHashString())
	}

	p.transactionCH <- tx

	return nil
}

func (p *Pending) collision(tx *transaction.Transaction) bool {
	_, ok := p.transactions[tx.GetHashString()]

	return ok
}

func (p *Pending) commitTx(tx *transaction.Transaction) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.transactions[tx.GetHashString()] = tx.Serialize()
	log.Printf("\"COMMIT-TX(%s)\": %s\n", p.pendingID, p.transactions[tx.GetHashString()])
}

func (p *Pending) flush() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.transactions = make(map[string]string)
}
