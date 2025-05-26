package mempool

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/core/signal"
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	PENDING_CLOSED_BUFFER_SIZE = 64
)

type MemPool struct {
	BlockTime time.Duration
	MaxTxSize uint32

	mu       sync.RWMutex
	pendings map[types.Topic]*Pending

	pendedCh chan<- *Pended
	closeCh  chan *signal.PendingClosing
}

func NewMemPool(blockTime time.Duration, maxTxSize uint32) *MemPool {
	mp := &MemPool{
		BlockTime: blockTime,
		MaxTxSize: maxTxSize,
		pendings:  make(map[types.Topic]*Pending),
	}

	go mp.closedPendingCollector()

	return mp
}

func (mp *MemPool) SetChannel(pendedCh chan<- *Pended) {
	mp.closeCh = make(chan *signal.PendingClosing, PENDING_CLOSED_BUFFER_SIZE)
	mp.pendedCh = pendedCh
}

func (mp *MemPool) AddPending(pendingId types.Topic, pendingTime time.Duration) error {
	if mp.IsOpen(pendingId) {
		return fmt.Errorf("topic (%s) already opened pending", pendingId)
	}

	pn := NewPending()

	pn.SetLimitOptions(mp.MaxTxSize, mp.BlockTime)
	pn.SetPendingOptions(pendingId, pendingTime)
	pn.SetChannel(mp.pendedCh, mp.closeCh)

	mp.AllocatePending(pendingId, pn)

	log.Printf(util.MemPoolString("MEMPOOL: New pending { topic: %s, duration: %s }"),
		pn.pendingID, pn.pendingTime)

	go pn.Activate()

	return nil
}

func (mp *MemPool) getPendingWithoutOpenCheck(pendingId types.Topic) *Pending {
	return mp.pendings[pendingId]
}

// Check Pending is opened
func (mp *MemPool) IsOpen(pendingId types.Topic) bool {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	_, ok := mp.pendings[pendingId]

	return ok
}

func (mp *MemPool) AllocatePending(pendingId types.Topic, open *Pending) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.pendings[pendingId] = open
}

func (mp *MemPool) CommitTransaction(pendingId types.Topic, tx *transaction.Transaction) error {
	if !mp.IsOpen(pendingId) {
		return fmt.Errorf("pending(%s) does not opened", pendingId)
	}

	pn := mp.getPendingWithoutOpenCheck(pendingId)

	if err := pn.PushTx(tx); err != nil {
		return err
	}

	return nil
}

func (mp *MemPool) closedPendingCollector() {
	for {
		log.Println(util.MemPoolString("MEMPOOL: ClosedPendingCollector blocked closeSignal waiting..."))

		c := <-mp.closeCh

		log.Printf(util.MemPoolString("MEMPOOL: ClosedPendingCollector closeSignal received: %s"), c.GetTopic())

		mp.closePending(c.GetTopic())

		c.Done()
	}
}

func (mp *MemPool) closePending(pendingId types.Topic) {
	mp.mu.Lock()
	log.Printf(util.MemPoolString("MEMPOOL: ClosedPendingCollector::closePending mutex locked: %s"), pendingId)

	defer mp.mu.Unlock()
	defer log.Printf(util.MemPoolString("MEMPOOL: ClosedPendingCollector::closePending mutex unlocked: %s"), pendingId)

	delete(mp.pendings, pendingId)
	log.Printf(util.MemPoolString("MEMPOOL: Pending %d remained"), len(mp.pendings))
}
