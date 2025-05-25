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
	// closeCh chan types.Topic
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
	// mp.closeCh = make(chan types.Topic, PENDING_CLOSED_BUFFER_SIZE)
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

	log.Printf(util.PendingString("Pending: New pending { topic: %s, duration: %s }"),
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
		log.Println(util.SystemString("ClosedPendingCollector looping"))

		c := <-mp.closeCh

		mp.mu.Lock()
		delete(mp.pendings, c.GetTopic())
		mp.mu.Unlock()

		c.Done()
	}
}
