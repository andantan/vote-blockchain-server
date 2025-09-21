package mempool

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type pendingMetaData struct {
	expired       bool
	cachedLength  int
	cachedOptions map[string]int
}

func newPendingMetaData(expired bool, cachedLength int, cachedOption map[string]int) *pendingMetaData {
	return &pendingMetaData{
		expired:       expired,
		cachedLength:  cachedLength,
		cachedOptions: cachedOption,
	}
}

type Pended struct {
	*pendingMetaData
	pendingID       types.Proposal
	pendingProposer types.Hash
	txx             map[string]*transaction.Transaction
}

func NewPended(pendingID types.Proposal, proposer types.Hash, txx map[string]*transaction.Transaction) *Pended {
	return &Pended{
		pendingMetaData: newPendingMetaData(false, 0, nil),
		pendingID:       pendingID,
		pendingProposer: proposer,
		txx:             txx,
	}
}

func NewExpiredPended(
	pendingID types.Proposal, cachedLength int, cachedOption map[string]int,
) *Pended {
	return &Pended{
		pendingMetaData: newPendingMetaData(true, cachedLength, cachedOption),
		pendingID:       pendingID,
	}
}

func (p *Pended) GetPendingID() types.Proposal {
	return p.pendingID
}

func (p *Pended) GetPendingProposer() types.Hash {
	return p.pendingProposer
}

func (p *Pended) GetTxx() map[string]*transaction.Transaction {
	return p.txx
}

func (p *Pended) GetCachedLength() int {
	return p.cachedLength
}

func (p *Pended) GetCachedOptions() map[string]int {
	return p.cachedOptions
}

func (p *Pended) IsExpired() bool {
	return p.expired
}
