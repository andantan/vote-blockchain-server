package mempool

import (
	"time"

	"github.com/andantan/vote-blockchain-server/types"
)

type PendingOpts struct {
	maxTxSize       uint32
	blockTime       time.Duration
	pendingID       types.Proposal
	pendingProposer types.Hash
	pendingTime     time.Duration
	pendedCh        chan *Pended
}

func NewPendingOpts(
	maxTxSize uint32,
	blockTime time.Duration,
	pendingID types.Proposal,
	proposer types.Hash,
	pendingTime time.Duration,
	pendedCh chan *Pended,
) *PendingOpts {
	return &PendingOpts{
		maxTxSize:       maxTxSize,
		blockTime:       blockTime,
		pendingID:       pendingID,
		pendingProposer: proposer,
		pendingTime:     pendingTime,
		pendedCh:        pendedCh,
	}
}
