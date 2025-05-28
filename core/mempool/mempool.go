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
		BlockTime: blockTime,
		MaxTxSize: maxTxSize,
		wg:        &sync.WaitGroup{},
		pendings:  make(map[types.Topic]*Pending),
	}

	mp.wg.Add(1)

	go mp.closedPendingCollector()

	return mp
}

func (mp *MemPool) SetChannel() {
	log.Printf(
		util.SystemString("SYSTEM: Memory pool setting channel... | { PENDED_REQUEST_BUFFER_SIZE: %d }"),
		PENDED_REQUEST_BUFFER_SIZE,
	)

	mp.pendedCh = make(chan *Pended, PENDED_REQUEST_BUFFER_SIZE)
	log.Println(util.SystemString("SYSTEM: Memory pool pended channel setting is done."))

	mp.shutdownCh = make(chan struct{})
	log.Println(util.SystemString("SYSTEM: Memory pool shutdown channel setting is done."))
}

func (mp *MemPool) Produce() <-chan *Pended {
	return mp.pendedCh
}

func (mp *MemPool) AddPending(pendingId types.Topic, pendingTime time.Duration) error {
	if mp.IsOpen(pendingId) {
		return fmt.Errorf("topic (%s) already opened pending", pendingId)
	}

	pn := NewPending()

	pn.SetLimitOptions(mp.MaxTxSize, mp.BlockTime)
	pn.SetPendingOptions(pendingId, pendingTime)
	pn.SetChannel(mp.pendedCh)

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
				if pending.IsTimeout() && !pending.IsClosed() {
					topicsOver = append(topicsOver, topic)
				} else if pending.IsClosed() {
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
