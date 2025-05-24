package mempool

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type Pended struct {
	pendingID types.Topic
	txx       map[string]*transaction.Transaction
}

func NewPended(pendingID types.Topic, txx map[string]*transaction.Transaction) *Pended {
	return &Pended{
		pendingID: pendingID,
		txx:       txx,
	}
}

func (p *Pended) GetPendingID() types.Topic {
	return p.pendingID
}

func (p *Pended) GetTxx() map[string]*transaction.Transaction {
	return p.txx
}
