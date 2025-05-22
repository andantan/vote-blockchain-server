package mempool

import (
	"fmt"
	"log"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type MemPool struct {
	BlockTime time.Duration
	MaxTxSize uint32
	pendings  map[types.Topic]*Pending
	closeCh   chan types.Topic
}

func NewMemPool(blockTime time.Duration, maxTxSize uint32) *MemPool {
	p := &MemPool{
		BlockTime: blockTime,
		MaxTxSize: maxTxSize,
		pendings:  make(map[types.Topic]*Pending),
		closeCh:   make(chan types.Topic),
	}

	go p.closedPendingCollector()

	return p
}

func (p *MemPool) AddPending(pendingId types.Topic, pendingTime time.Duration) error {
	if p.IsOpen(pendingId) {
		return fmt.Errorf("topic(%s) already opened pending", pendingId)
	}

	pn := NewPending(pendingTime, p.BlockTime, p.MaxTxSize, pendingId, p.closeCh)

	p.AllocatePending(pendingId, pn)

	log.Printf("topic(%s) pending success, duration (%s)", pn.pendingID, pn.pendingTime)

	go pn.Activate()

	return nil
}

func (p *MemPool) getPendingWithoutOpenCheck(pendingId types.Topic) *Pending {
	return p.pendings[pendingId]
}

// Check Pending is opened
func (p *MemPool) IsOpen(pendingId types.Topic) bool {
	_, ok := p.pendings[pendingId]

	return ok
}

func (p *MemPool) AllocatePending(pendingId types.Topic, open *Pending) {
	p.pendings[pendingId] = open
}

func (p *MemPool) CommitTransaction(pendingId types.Topic, tx *transaction.Transaction) error {
	if !p.IsOpen(pendingId) {
		return fmt.Errorf("pending(%s) does not opened", pendingId)
	}

	pn := p.getPendingWithoutOpenCheck(pendingId)

	if err := pn.CommitTx(tx); err != nil {
		return err
	}

	return nil
}

func (p *MemPool) closedPendingCollector() {
	for {
		topic := <-p.closeCh
		delete(p.pendings, topic)

		log.Printf("pending (%s) removed from memPool\n", topic)
	}
}
