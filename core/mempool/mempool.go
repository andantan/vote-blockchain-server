package mempool

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	PENDED_REQUEST_BUFFER_SIZE = 128
)

type MemPool struct {
	BlockTime time.Duration
	MaxTxSize uint32

	wg       *sync.WaitGroup
	mu       sync.RWMutex
	pendings map[types.Topic]*Pending

	pendedCh   chan *Pended
	shutdownCh chan struct{}
}

func NewMemPool(blockTime time.Duration, maxTxSize uint32) *MemPool {
	mp := &MemPool{
		BlockTime:  blockTime,
		MaxTxSize:  maxTxSize,
		wg:         &sync.WaitGroup{},
		pendings:   make(map[types.Topic]*Pending),
		pendedCh:   make(chan *Pended, PENDED_REQUEST_BUFFER_SIZE),
		shutdownCh: make(chan struct{}),
	}

	mp.wg.Add(1)

	go mp.closedPendingCollector()

	return mp
}

func (mp *MemPool) Consume() <-chan *Pended {
	return mp.pendedCh
}

func (mp *MemPool) AddPending(pendingId types.Topic, pendingTime time.Duration) error {
	if mp.IsOpen(pendingId) {
		return fmt.Errorf("topic (%s) already opened pending", pendingId)
	}

	pnOpts := NewPendingOpts(mp.MaxTxSize, mp.BlockTime, pendingId, pendingTime, mp.pendedCh)
	pn := NewPending(pnOpts)

	mp.openPending(pendingId, pn)

	go pn.Activate()

	log.Printf(util.MemPoolString("MEMPOOL: New pending { topic: %s, duration: %s }"),
		pn.pendingID, pn.pendingTime)

	return nil
}

func (mp *MemPool) getPendingWithoutOpenCheck(pendingId types.Topic) *Pending {
	return mp.pendings[pendingId]
}

func (mp *MemPool) IsOpen(pendingId types.Topic) bool {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	_, ok := mp.pendings[pendingId]

	return ok
}

func (mp *MemPool) openPending(pendingId types.Topic, open *Pending) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.pendings[pendingId] = open
}

func (mp *MemPool) CommitTransaction(pendingId types.Topic, tx *transaction.Transaction) error {
	if !mp.IsOpen(pendingId) {
		return fmt.Errorf("pending(%s) does not opened", pendingId)
	}

	pn := mp.getPendingWithoutOpenCheck(pendingId)

	if pn.IsTimeout() {
		return fmt.Errorf("pending(%s) time over", pendingId)
	}

	if pn.IsClosed() {
		return fmt.Errorf("pending(%s) is closed", pendingId)
	}

	if err := pn.PushTx(tx); err != nil {
		return err
	}

	return nil
}

func (mp *MemPool) closedPendingCollector() {
	defer mp.wg.Done()

	collector := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-collector.C:

			topicsAlive := []types.Topic{}
			topicsOver := []types.Topic{}
			topicsToRemove := []types.Topic{}

			for topic, pending := range mp.pendings {
				isTimeout := pending.IsTimeout()
				isClosed := pending.IsClosed()

				if isTimeout && !isClosed {
					topicsOver = append(topicsOver, topic)
				} else if isClosed {
					topicsToRemove = append(topicsToRemove, topic)
				} else {
					topicsAlive = append(topicsAlive, topic)
				}
			}

			log.Printf(
				util.MemPoolString("MEMPOOL: Collector scanned. Alive: %d, Timeout: %d, To Remove: %d"),
				len(topicsAlive),
				len(topicsOver),
				len(topicsToRemove),
			)

			for _, topic := range topicsToRemove {
				mp.closePending(topic)
			}
		case <-mp.shutdownCh:
			for _, pending := range mp.pendings {
				pending.ctx.Done()
			}

			return
		}
	}
}

func (mp *MemPool) closePending(pendingId types.Topic) {
	mp.mu.Lock()
	delete(mp.pendings, pendingId)
	log.Printf(util.MemPoolString("MEMPOOL: %s| Pending successfully removed"), pendingId)
	mp.mu.Unlock()
}

func (mp *MemPool) Shutdown() {
	log.Println(util.MemPoolString("MEMPOOL: Initiating shutdown for MemPool"))
	close(mp.shutdownCh)
	mp.wg.Wait()
	close(mp.pendedCh)
	log.Println(util.MemPoolString("MEMPOOL: MemPool shutdown complete"))
}
