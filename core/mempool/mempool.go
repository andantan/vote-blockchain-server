package mempool

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/config"
	"github.com/andantan/vote-blockchain-server/core/transaction"
	werror "github.com/andantan/vote-blockchain-server/error"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
)

type MemPool struct {
	BlockTime time.Duration
	MaxTxSize uint32

	wg       *sync.WaitGroup
	mu       sync.RWMutex
	pendings map[types.Proposal]*Pending

	pendedCh   chan *Pended
	shutdownCh chan struct{}
}

func NewMemPool() *MemPool {
	__cfg := config.GetChainParameterConfiguration()
	__sys_channel_size := config.GetChannelBufferSizeSystemConfiguration()

	blockInterval := time.Duration(__cfg.BlockIntervalSeconds) * time.Second

	mp := &MemPool{
		BlockTime: blockInterval,
		MaxTxSize: __cfg.MaxTransactionSize,
		wg:        &sync.WaitGroup{},
		pendings:  make(map[types.Proposal]*Pending),
		pendedCh: make(
			chan *Pended,
			__sys_channel_size.PendedPropaginateChannelBufferSize,
		),
		shutdownCh: make(chan struct{}),
	}

	mp.wg.Add(1)

	go mp.closedPendingCollector()

	return mp
}

func (mp *MemPool) Consume() <-chan *Pended {
	return mp.pendedCh
}

func (mp *MemPool) AddPending(pendingId types.Proposal, pendingTime time.Duration) error {
	if mp.IsOpen(pendingId) {
		msg := fmt.Sprintf("Given proposal (%s) is already pending", pendingId)
		return werror.NewWrappedError("PROPOSAL_ALREADY_OPEN", msg, nil)
	}

	pnOpts := NewPendingOpts(mp.MaxTxSize, mp.BlockTime, pendingId, pendingTime, mp.pendedCh)
	pn := NewPending(pnOpts)

	mp.openPending(pendingId, pn)

	go pn.Activate()

	log.Printf(util.MemPoolString("MEMPOOL: New pending { topic: %s, duration: %s }"),
		pn.pendingID, pn.pendingTime)

	return nil
}

func (mp *MemPool) getPendingWithoutOpenCheck(pendingId types.Proposal) *Pending {
	return mp.pendings[pendingId]
}

func (mp *MemPool) IsOpen(pendingId types.Proposal) bool {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	_, ok := mp.pendings[pendingId]

	return ok
}

func (mp *MemPool) openPending(pendingId types.Proposal, open *Pending) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.pendings[pendingId] = open
}

func (mp *MemPool) CommitTransaction(pendingId types.Proposal, tx *transaction.Transaction) error {
	if !mp.IsOpen(pendingId) {
		msg := fmt.Sprintf("Proposal (%s) is not open", pendingId)
		return werror.NewWrappedError("PROPOSAL_NOT_OPEN", msg, nil)
	}

	pn := mp.getPendingWithoutOpenCheck(pendingId)

	if pn.IsClosed() {
		msg := fmt.Sprintf("Proposal (%s) is closed", pendingId)
		return werror.NewWrappedError("CLOSED_PROPOSAL", msg, nil)
	}

	if pn.IsTimeout() {
		msg := fmt.Sprintf("Proposal (%s) has timed out", pendingId)
		return werror.NewWrappedError("TIMEOUT_PROPOSAL", msg, nil)
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

			topicsAlive := []types.Proposal{}
			topicsOver := []types.Proposal{}
			topicsToRemove := []types.Proposal{}

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
			return
		}
	}
}

func (mp *MemPool) closePending(pendingId types.Proposal) {
	mp.mu.Lock()
	delete(mp.pendings, pendingId)
	log.Printf(util.MemPoolString("MEMPOOL: %s| Pending successfully removed"), pendingId)
	mp.mu.Unlock()
}

func (mp *MemPool) shutDownclosedPendingCollector() {
	close(mp.shutdownCh)
	mp.wg.Wait()
}

func (mp *MemPool) Shutdown() {
	log.Println(util.MemPoolString("MEMPOOL: Initiating shutdown for MemPool"))
	mp.shutDownclosedPendingCollector()

	mwg := &sync.WaitGroup{}
	mwg.Add(len(mp.pendings))

	for _, pending := range mp.pendings {
		go pending.Shutdown(mwg)
	}

	mwg.Wait()

	log.Println(util.MemPoolString("MEMPOOL: All of pendings shutdown complete"))

	close(mp.pendedCh)

	log.Println(util.MemPoolString("MEMPOOL: MemPool shutdown complete"))
}
