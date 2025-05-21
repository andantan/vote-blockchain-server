package mempool

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type MemPool interface {
}

type PendingPool struct {
	pendings map[types.VotingID]*Pending
}

func NewPendingPool() *PendingPool {
	return &PendingPool{
		pendings: make(map[types.VotingID]*Pending),
	}
}

func (p *PendingPool) AddPending(tx *transaction.Transaction, id types.VotingID) error {

	return nil
}
