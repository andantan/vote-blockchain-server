package mempool

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type pendingMetaData struct {
	expired       bool
	cachedOptions map[string]int
}

func newPendingMetaData(expire bool, cachedOption map[string]int) *pendingMetaData {
	return &pendingMetaData{
		expired:       expire,
		cachedOptions: cachedOption,
	}
}

type Pended struct {
	*pendingMetaData
	pendingID types.Topic
	txx       map[string]*transaction.Transaction
}

func NewPended(pendingID types.Topic, txx map[string]*transaction.Transaction) *Pended {
	return &Pended{
		pendingID:       pendingID,
		txx:             txx,
		pendingMetaData: newPendingMetaData(false, nil),
	}
}

func NewExpiredPended(
	pendingID types.Topic, txx map[string]*transaction.Transaction,
	cachedOption map[string]int,
) *Pended {
	return &Pended{
		pendingID:       pendingID,
		txx:             txx,
		pendingMetaData: newPendingMetaData(true, cachedOption),
	}
}

func (p *Pended) GetPendingID() types.Topic {
	return p.pendingID
}

func (p *Pended) GetTxx() map[string]*transaction.Transaction {
	return p.txx
}

func (p *Pended) GetCachedOptions() map[string]int {
	return p.cachedOptions
}

func (p *Pended) IsExpired() bool {
	return p.expired
}
