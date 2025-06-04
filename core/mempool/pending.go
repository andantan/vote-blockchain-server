package mempool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	werror "github.com/andantan/vote-blockchain-server/error"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	RESET_TIMER_DURATION = iota
)

const (
	Tx_BUFFER_SIZE = 1024
)

const (
	INTERRUPT_TIMER_DURATION = 5 * time.Second
	CLOSE_TIMER_DURATION     = 10 * time.Second
)

type Pending struct {
	pendingID types.Proposal // TopicID

	mu       sync.RWMutex
	txx      map[string]*transaction.Transaction
	txCache  map[string]struct{}
	optCache map[string]int

	timeout bool // txCh flag
	closed  bool

	wg       *sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	ctxTimer *time.Timer

	blockTime   time.Duration // Block Time (system)
	blockTicker *time.Ticker

	pendingTime   time.Duration // Vote duration
	pendingTicker *time.Ticker

	closeTimer *time.Timer

	maxTxSize uint32 // Tx size (system)

	transactionCh chan *transaction.Transaction
	pendedCh      chan *Pended
}

func NewPending(opts *PendingOpts) *Pending {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pending{
		pendingID:     opts.pendingID,
		txx:           make(map[string]*transaction.Transaction),
		txCache:       make(map[string]struct{}),
		optCache:      make(map[string]int),
		timeout:       false,
		closed:        false,
		wg:            &sync.WaitGroup{},
		ctx:           ctx,
		cancel:        cancel,
		pendingTime:   opts.pendingTime,
		blockTime:     opts.blockTime,
		maxTxSize:     opts.maxTxSize,
		transactionCh: make(chan *transaction.Transaction, Tx_BUFFER_SIZE),
		pendedCh:      opts.pendedCh,
	}

	p.wg.Add(1)

	return p
}

func (p *Pending) Len() int {
	p.mu.Lock()
	l := len(p.txx)
	p.mu.Unlock()

	return l
}

func (p *Pending) Activate() {
	defer func() {
		p.wg.Done()
		log.Printf(util.PendingString("PENDING: %s | Activation exited"), p.pendingID)
	}()

	p.blockTicker = time.NewTicker(p.blockTime)
	p.pendingTicker = time.NewTicker(p.pendingTime)

	p.closeTimer = time.NewTimer(RESET_TIMER_DURATION)
	<-p.closeTimer.C
	p.ctxTimer = time.NewTimer(RESET_TIMER_DURATION)
	<-p.ctxTimer.C

	defer func() {
		p.blockTicker.Stop()
		p.pendingTicker.Stop()
		p.closeTimer.Stop()
		p.ctxTimer.Stop()
	}()

	for {
		select {
		case tx := <-p.transactionCh:
			p.processTx(tx, true)

		case <-p.blockTicker.C:
			p.flushIfNotEmpty()

		case <-p.pendingTicker.C:
			p.timeoutPending()
			p.stopBlockTicker()
			p.startCloseTimer()

		case <-p.closeTimer.C:
			log.Printf(util.PendingString("PENDING: %s | Pending is now on closing"), p.pendingID)

			p.closeTxChannel()
			p.processClosedTxCh()
			p.flushIfNotEmpty()
			p.emitExpiredPended()
			p.clearCache()
			p.closePending()

			return

		case <-p.ctx.Done():
			log.Printf(
				util.PendingString("PENDING: %s | Pending interrupted. Initiating immediate shutdown due to external signal."),
				p.pendingID,
			)

			p.interruptPending()
			p.stopBlockTicker()
			p.stopPendingTicker()
			p.stopCloseTimer()
			p.startContextTimer()

			<-p.ctxTimer.C

			p.closeTxChannel()
			p.processClosedTxCh()
			p.flushIfNotEmpty()
			p.emitExpiredPended()
			p.clearCache()
			p.closePending()

			return
		}
	}
}

func (p *Pending) Shutdown(mwg *sync.WaitGroup) {
	defer mwg.Done()

	log.Printf(util.PendingString("PENDING: %s | service shutdown initiated"), p.pendingID)

	p.cancel()

	p.wg.Wait()

	log.Printf(util.PendingString("PENDING: %s | service shutdown completed"), p.pendingID)
}

func (p *Pending) PushTx(tx *transaction.Transaction) error {
	if p.collision(tx) {
		msg := fmt.Sprintf("Vote transaction with hash '%s' for proposal '%s' has already been submitted.", tx.GetHashString(), p.pendingID)
		return werror.NewWrappedError("DUPLICATE_VOTE_SUBMISSION", msg, nil)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case p.transactionCh <- tx:
		return nil
	default:
		// TODO Retrive transaction
		return fmt.Errorf("failed to push tx %s | %s: transaction channel is likely full or closed during send attempt",
			p.pendingID,
			tx.GetHashString(),
		)
	}

}

func (p *Pending) processTx(tx *transaction.Transaction, blockTickerReset bool) {
	p.commitTx(tx)

	if p.maxTxSize <= uint32(p.Len()) {
		p.emitAndFlush()

		if blockTickerReset {
			p.blockTicker.Reset(p.blockTime)
		}
	}
}

func (p *Pending) closeTxChannel() {
	close(p.transactionCh)

	log.Printf(util.PendingString("PENDING: %s | Transaction channel closed."), p.pendingID)
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

func (p *Pending) emitExpiredPended() {
	p.pendedCh <- NewExpiredPended(p.pendingID, len(p.txCache), p.optCache)
}

func (p *Pending) flushIfNotEmpty() {
	if uint32(p.Len()) != 0 {
		p.emitAndFlush()
	}
}

func (p *Pending) processClosedTxCh() {
	log.Printf(util.PendingString("PENDING: %s | Flushing closed tx channel"), p.pendingID)

	for tx := range p.transactionCh {
		p.processTx(tx, false)

		log.Printf(util.PendingString("PENDING: %s | Flushed remaining transaction during shutdown"), p.pendingID)
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
	p.cacheOption(tx)
}

func (p *Pending) commit(tx *transaction.Transaction) {
	p.txx[tx.GetHashString()] = tx
}

func (p *Pending) cache(tx *transaction.Transaction) {
	p.txCache[tx.GetHashString()] = struct{}{}
}

func (p *Pending) cacheOption(tx *transaction.Transaction) {
	p.optCache[tx.Option]++
}

func (p *Pending) clearCache() {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf(util.PendingString("PENDING: %s | Cache { txCachedLength: %d, txCachedOption: %v }"), p.pendingID, len(p.txCache), p.optCache)

	p.txCache = make(map[string]struct{})
	p.optCache = make(map[string]int)
}

func (p *Pending) timeoutPending() {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf(util.PendingString("PENDING: %s | Pending is over"), p.pendingID)

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

func (p *Pending) startCloseTimer() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closeTimer.Reset(CLOSE_TIMER_DURATION)
}

func (p *Pending) startContextTimer() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.ctxTimer.Reset(INTERRUPT_TIMER_DURATION)
}

func (p *Pending) stopBlockTicker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.blockTicker.Stop()
}

func (p *Pending) stopPendingTicker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pendingTicker.Stop()
}

func (p *Pending) stopCloseTimer() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closeTimer.Stop()
}
