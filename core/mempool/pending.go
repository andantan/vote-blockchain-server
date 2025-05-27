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

type Pending struct {
	pendingID types.Topic // TopicID

	mu          sync.RWMutex
	txx         map[string]*transaction.Transaction
	txCache     map[string]struct{}
	txChannling bool // txCh flag

	pendingTime time.Duration // Vote duration
	blockTime   time.Duration // Block Time (system)
	maxTxSize   uint32        // Tx size (system)
	//scheduledBlockHeight []uint64      // Pended block heights

	transactionCH  chan *transaction.Transaction
	pendedCh       chan<- *Pended
	pendingCloseCh chan<- *signal.PendingClosing
}

func NewPending() *Pending {
	p := &Pending{
		txx:     make(map[string]*transaction.Transaction),
		txCache: make(map[string]struct{}),
	}

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
	pendingCloseCh chan<- *signal.PendingClosing,
) {
	p.transactionCH = make(chan *transaction.Transaction)
	p.txChannling = true

	p.pendedCh = pendedCh
	p.pendingCloseCh = pendingCloseCh
}

func (p *Pending) Len() int {
	p.mu.Lock()
	l := len(p.txx)
	p.mu.Unlock()

	return l
}

func (p *Pending) Activate() {
	blockTimer := time.NewTicker(p.blockTime)
	pendingTimer := time.NewTicker(p.pendingTime)

	defer log.Printf(util.PendingString("Pending: %s | Activation exited"), p.pendingID)
	defer blockTimer.Stop()
	defer pendingTimer.Stop()

	for {
		select {
		case tx, ok := <-p.transactionCH:
			if !ok {
				p.stopReceivingTx(false)
				p.closedTxChFlush()
				p.clearTxCache()
				return
			}

			p.processTxWithResetTimer(tx, blockTimer)

		case <-blockTimer.C:
			p.flushIfNotEmpty()

		case <-pendingTimer.C:
			log.Printf(util.PendingString("Pending: %s | Pending is over"), p.pendingID)

			p.stopReceivingTx(true)
			p.closedTxChFlush()
			p.flushIfNotEmpty()
			p.clearTxCache()
			return
		}
	}
}

func (p *Pending) PushTx(tx *transaction.Transaction) error {
	if p.collision(tx) {
		return fmt.Errorf("given tx (%s) is already commited", tx.GetHashString())
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.txChannling {
		return fmt.Errorf("pending %s is not channeling transactions (closed or shutting down)", p.pendingID)
	}

	select {
	case p.transactionCH <- tx:
		return nil
	default:
		return fmt.Errorf("failed to push tx %s | %s: transaction channel is likely full or closed during send attempt",
			p.pendingID,
			tx.GetHashString(),
		)
	}

}

func (p *Pending) processTxWithResetTimer(tx *transaction.Transaction, blockTimer *time.Ticker) {
	p.commitTx(tx)

	if p.maxTxSize <= uint32(p.Len()) {
		p.emitAndFlush()
		blockTimer.Reset(p.blockTime)
	}
}

func (p *Pending) stopReceivingTx(txChClose bool) {
	if txChClose {
		close(p.transactionCH)

		log.Printf(util.PendingString("Pending: %s | Transaction channel closed."), p.pendingID)
	}

	p.unChanneling()
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

func (p *Pending) flushIfNotEmpty() {
	if uint32(p.Len()) != 0 {
		p.emitAndFlush()
	}
}

func (p *Pending) closedTxChFlush() {
	log.Printf(
		util.PendingString("Pending: %s | Transaction channel closed by sender. Exiting transaction processing loop."),
		p.pendingID,
	)

	for tx := range p.transactionCH {
		p.commitTx(tx)

		if p.maxTxSize <= uint32(p.Len()) {
			p.emitAndFlush()
		}

		log.Printf(util.PendingString("Pending: %s | Flushed remaining transaction during shutdown."), p.pendingID)
	}
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
}

func (p *Pending) commit(tx *transaction.Transaction) {
	p.txx[tx.GetHashString()] = tx
}

func (p *Pending) cache(tx *transaction.Transaction) {
	p.txCache[tx.GetHashString()] = struct{}{}
}

func (p *Pending) clearTxCache() {
	p.mu.Lock()
	p.txCache = make(map[string]struct{})
	defer p.mu.Unlock()

	log.Printf(util.PendingString("Pending: %s | TxCache cleared for memory optimization."), p.pendingID)
}

func (p *Pending) IsPendingChanneled() bool {
	p.mu.RLock()
	c := p.txChannling
	p.mu.RUnlock()

	return c
}

func (p *Pending) unChanneling() {
	p.mu.Lock()
	p.txChannling = false
	p.mu.Unlock()
}
