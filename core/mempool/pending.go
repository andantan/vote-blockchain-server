package mempool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/core/signal"
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
)

type Pending struct {
	pendingID types.Topic // TopicID

	mu      sync.RWMutex
	txx     map[string]*transaction.Transaction
	txCache map[string]struct{}

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	pendingTime          time.Duration // Vote duration
	blockTime            time.Duration // Block Time (system)
	maxTxSize            uint32        // Tx size (system)
	scheduledBlockHeight []uint64      // Pended block heights

	transactionCH chan *transaction.Transaction
	pendedCh      chan<- *Pended
	closeCh       chan<- *signal.PendingClosing
	// closeCh chan<- types.Topic
}

func NewPending() *Pending {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pending{
		txx:                  make(map[string]*transaction.Transaction),
		txCache:              make(map[string]struct{}),
		scheduledBlockHeight: []uint64{},
		ctx:                  ctx,
		cancel:               cancel,
		wg:                   &sync.WaitGroup{},
		transactionCH:        make(chan *transaction.Transaction),
	}

	p.wg.Add(1)

	return p
}

func (p *Pending) SetLimitOptions(maxTxSize uint32, blockTime time.Duration) {
	p.maxTxSize = maxTxSize
	p.blockTime = blockTime
}

func (p *Pending) SetPendingOptions(pendingId types.Topic, pendingTime time.Duration) {
	p.pendingID = pendingId
	p.pendingTime = pendingTime
}

func (p *Pending) SetChannel(
	pendedCh chan<- *Pended,
	closeCh chan<- *signal.PendingClosing,
	// closeCh chan<- types.Topic,
) {
	p.pendedCh = pendedCh
	p.closeCh = closeCh
}

func (p *Pending) Len() int {
	p.mu.Lock()
	l := len(p.txx)
	p.mu.Unlock()

	return l
}

func (p *Pending) Transactions() *transaction.SortedTxx {
	s := transaction.NewSortedTxx(p.txx)

	return s
}

func (p *Pending) Activate() {
	defer p.wg.Done()

	blockTimer := time.NewTicker(p.blockTime)
	defer blockTimer.Stop()
	pendingTimer := time.NewTicker(p.pendingTime)
	defer pendingTimer.Stop()

	for {
		select {
		case <-p.ctx.Done():
			log.Printf(util.PendingString("Pending: %s | Activation exited."), p.pendingID)

			return

		case tx := <-p.transactionCH:
			p.commitTx(tx)

			if p.maxTxSize <= uint32(p.Len()) {
				// log.Printf(util.PendingString("PENDING: %s | Max transactions (%d) reached"),
				// 	p.pendingID, p.Len())

				p.emitAndFlush()

				// log.Printf(util.BlockString("Block: %s | Reached maxTxSize - New block created"), p.pendingID)
				// log.Printf(util.PendingString("PENDING: %s | Transactions map cleared. New size: %d"),
				// 	p.pendingID, p.Len())

				blockTimer.Reset(p.blockTime)
			}

		case <-blockTimer.C:
			if uint32(p.Len()) != 0 {
				p.emitAndFlush()
				// log.Printf(util.BlockString("Block: %s | Block timeout - new block created"), p.pendingID)
			}
		case <-pendingTimer.C:
			if uint32(p.Len()) != 0 {
				p.emitAndFlush()
			}

			//log.Printf(util.PendingString("Pending: %s | Pending is over"), p.pendingID)

			go p.triggerShutdown()

			return
		}
	}

}

func (p *Pending) PushTx(tx *transaction.Transaction) error {
	if p.collision(tx) {
		return fmt.Errorf("given tx (%s) is already commited", tx.GetHashString())
	}

	p.transactionCH <- tx

	return nil
}

func (p *Pending) seekTx(hash string) *transaction.Transaction {
	p.mu.RLock()
	defer p.mu.RUnlock()
	t, ok := p.txx[hash]
	if !ok {
		return nil
	}

	s := transaction.NewTransaction(
		t.GetHash(),
		t.GetOption(),
		t.GetTimeStamp(),
	)

	return s
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

	// log.Printf(util.CommitString("COMMIT: %s | New commit tx { hash: %s, option: %s, timestamp: %d }"),
	// 	p.pendingID, tx.GetHashString(), tx.GetOption(), tx.GetTimeStamp())
}

func (p *Pending) commit(tx *transaction.Transaction) {
	p.txx[tx.GetHashString()] = tx
}

func (p *Pending) cache(tx *transaction.Transaction) {
	p.txCache[tx.GetHashString()] = struct{}{}
}

func (p *Pending) flush() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.txx = make(map[string]*transaction.Transaction)
}

func (p *Pending) emitAndFlush() {
	p.pendedCh <- NewPended(p.pendingID, p.txx)
	p.flush()
}

func (p *Pending) triggerShutdown() {
	p.cancel()
	p.wg.Wait()

	s := signal.NewPendingClosing(p.pendingID)

	s.Add(1)
	p.closeCh <- s
	s.Wait()

	log.Printf(util.PendingString("Pending: %s | successfully closed and removed from MemPool"), p.pendingID)
}
