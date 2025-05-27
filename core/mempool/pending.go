package mempool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	RESET_TIMER_DURATION = iota
)

const (
	CONTEXT_TIMER_DURATION = 5 * time.Second
	CLOSE_TIMER_DURATION   = 10 * time.Second
)

type Pending struct {
	pendingID types.Topic // TopicID

	mu      sync.RWMutex
	txx     map[string]*transaction.Transaction
	txCache map[string]struct{}

	timeout bool // txCh flag
	closed  bool

	ctx    context.Context
	cancel context.CancelFunc

	pendingTime time.Duration // Vote duration
	blockTime   time.Duration // Block Time (system)
	maxTxSize   uint32        // Tx size (system)
	//scheduledBlockHeight []uint64      // Pended block heights

	transactionCH chan *transaction.Transaction
	pendedCh      chan<- *Pended
}

func NewPending() *Pending {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pending{
		txx:     make(map[string]*transaction.Transaction),
		txCache: make(map[string]struct{}),
		timeout: false,
		closed:  false,
		ctx:     ctx,
		cancel:  cancel,
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

func (p *Pending) SetChannel(pendedCh chan<- *Pended) {
	p.transactionCH = make(chan *transaction.Transaction)

	p.pendedCh = pendedCh
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
	closeTimer := time.NewTimer(RESET_TIMER_DURATION)
	<-closeTimer.C
	contextTimer := time.NewTimer(RESET_TIMER_DURATION)
	<-contextTimer.C

	defer log.Printf(util.PendingString("Pending: %s | Activation exited"), p.pendingID)
	defer blockTimer.Stop()
	defer pendingTimer.Stop()
	defer closeTimer.Stop()
	defer contextTimer.Stop()

	for {
		select {
		case <-p.ctx.Done():
			log.Printf(
				util.PendingString("Pending: %s | Context cancelled. Initiating immediate shutdown due to external signal."),
				p.pendingID,
			)
			blockTimer.Stop()
			pendingTimer.Stop()
			closeTimer.Stop()

			p.interruptPending()

			contextTimer.Reset(CONTEXT_TIMER_DURATION)

		case <-contextTimer.C:
			log.Printf(
				util.PendingString("Pending: %s | Grace period (%s) for transaction finalization ended. Closing transaction channel to prevent panics."),
				p.pendingID,
				CONTEXT_TIMER_DURATION,
			)

			p.closeTxChannel()
			p.closedTxChFlush()
			p.flushIfNotEmpty()
			p.clearTxCache()

			return

		case tx, ok := <-p.transactionCH:
			if !ok {
				p.timeoutPending()
				p.closedTxChFlush()
				p.clearTxCache()
				p.closePending()
				return
			}

			p.processTxWithResetTimer(tx, blockTimer)

		case <-blockTimer.C:
			p.flushIfNotEmpty()

		case <-pendingTimer.C:
			log.Printf(util.PendingString("Pending: %s | Pending is over"), p.pendingID)

			p.timeoutPending()
			closeTimer.Reset(CONTEXT_TIMER_DURATION)

		case <-closeTimer.C:
			log.Printf(util.PendingString("Pending: %s | Pending is now on closing"), p.pendingID)

			p.closeTxChannel()
			p.closedTxChFlush()
			p.flushIfNotEmpty()
			p.clearTxCache()
			p.closePending()
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

func (p *Pending) closeTxChannel() {
	close(p.transactionCH)

	log.Printf(util.PendingString("Pending: %s | Transaction channel closed."), p.pendingID)
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
		util.PendingString("Pending: %s | Flushing closed tx channel"),
		p.pendingID,
	)

	for tx := range p.transactionCH {
		p.commitTx(tx)

		if p.maxTxSize <= uint32(p.Len()) {
			p.emitAndFlush()
		}

		log.Printf(util.PendingString("Pending: %s | Flushed remaining transaction during shutdown"), p.pendingID)
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
	defer p.mu.Unlock()

	txCachedLength := len(p.txCache)

	p.txCache = make(map[string]struct{})

	clearedTxCacheLength := len(p.txCache)

	log.Printf(util.PendingString("Pending: %s | TxCache cleared { txCachedLength: %d, clearedTxCacheLength: %d }"),
		p.pendingID, txCachedLength, clearedTxCacheLength)
}

func (p *Pending) IsTimeout() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.timeout
}

func (p *Pending) IsClosed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.closed
}

func (p *Pending) timeoutPending() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.timeout = true
}

func (p *Pending) closePending() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closed = true
}

func (p *Pending) interruptPending() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.timeout = true
	p.closed = true
}
