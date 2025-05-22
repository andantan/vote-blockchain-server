package mempool

import (
	"fmt"
	"log"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type Pending struct {
	transactions         map[types.Hash]string
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
		transactions:         make(map[types.Hash]string),
		transactionCH:        make(chan *transaction.Transaction),
		pendingTime:          pendingTime,
		blockTime:            blockTime,
		maxTransactionSize:   maxTxSize,
		scheduledBlockHeight: []uint64{},
		pendingID:            pendingId,
		closeCh:              closeCh,
	}
}

func (p *Pending) Activate() {
	blockTimer := time.NewTicker(p.blockTime)
	pendingTimer := time.NewTicker(p.pendingTime)

labelPending:
	for {
		select {
		case tx := <-p.transactionCH:
			log.Printf("New tx in pending(%s): %s|%s\n", p.pendingID, tx.Hash.String(), tx.Option)
		case <-blockTimer.C:
			log.Printf("New (%s) block created\n", p.pendingID)
		case <-pendingTimer.C:
			log.Printf("Pending(%s) over\n", p.pendingID)
			log.Printf("New block created(%s) by close pending\n", p.pendingID)

			p.closeCh <- p.pendingID

			break labelPending
		}
	}
}

func (p *Pending) CommitTx(tx *transaction.Transaction) error {
	if p.collision(tx) {
		return fmt.Errorf("given tx (%s) is already commited", tx.Hash.String())
	}

	p.transactionCH <- tx

	return nil
}

func (p *Pending) collision(tx *transaction.Transaction) bool {
	_, ok := p.transactions[tx.Hash]

	return ok
}
