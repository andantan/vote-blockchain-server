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

	pendedCh       chan<- *Pended
	pendingCloseCh chan *signal.PendingClosing
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
	log.Printf(
		util.SystemString("SYSTEM: Memory pool setting channel... | { PENDING_CLOSED_BUFFER_SIZE: %d }"),
		PENDING_CLOSED_BUFFER_SIZE,
	)

	mp.pendedCh = pendedCh
	log.Println(util.SystemString("SYSTEM: Memory pool pended channel setting is done."))
	mp.pendingCloseCh = make(chan *signal.PendingClosing, PENDING_CLOSED_BUFFER_SIZE)
	log.Println(util.SystemString("SYSTEM: Memory pool pendingClose channel setting is done."))
}

func (mp *MemPool) AddPending(pendingId types.Topic, pendingTime time.Duration) error {
	if mp.IsOpen(pendingId) {
		return fmt.Errorf("topic (%s) already opened pending", pendingId)
	}

	pn := NewPending()

	pn.SetLimitOptions(mp.MaxTxSize, mp.BlockTime)
	pn.SetPendingOptions(pendingId, pendingTime)
	pn.SetChannel(mp.pendedCh, mp.pendingCloseCh)

	mp.AllocatePending(pendingId, pn)

	go pn.Activate()

	log.Printf(util.MemPoolString("MEMPOOL: New pending { topic: %s, duration: %s }"),
		pn.pendingID, pn.pendingTime)

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

	if !pn.IsPendingChanneled() {
		return fmt.Errorf("pending(%s) time over", pendingId)
	}

	if err := pn.PushTx(tx); err != nil {
		return err
	}

	return nil
}

func (mp *MemPool) closedPendingCollector() {
	collector := time.NewTicker(5 * time.Second)

	for {
		<-collector.C

		topicsAlive := []types.Topic{}
		topicsToRemove := []types.Topic{}

		for topic, pending := range mp.pendings {
			if !pending.IsPendingChanneled() {
				// mp.closePending(topic)
				topicsToRemove = append(topicsToRemove, topic)
			} else {
				topicsAlive = append(topicsAlive, topic)
			}
		}

		log.Printf(
			util.MemPoolString("MEMPOOL: Collector scanned. Alive: %d, To Remove: %d"),
			len(topicsAlive),
			len(topicsToRemove),
		)

		for _, topic := range topicsToRemove {
			mp.closePending(topic)
		}
	}
}

func (mp *MemPool) closePending(pendingId types.Topic) {
	mp.mu.Lock()
	delete(mp.pendings, pendingId)
	log.Printf(util.MemPoolString("MEMPOOL: Pending %s successfully removed"), pendingId)
	mp.mu.Unlock()
}
