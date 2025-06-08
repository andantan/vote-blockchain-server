package mempool

import (
	"time"

	"github.com/andantan/vote-blockchain-server/types"
)

type PendingOpts struct {
	maxTxSize   uint32
	blockTime   time.Duration
	pendingID   types.Proposal
	pendingTime time.Duration
	pendedCh    chan *Pended
}

func NewPendingOpts(
	maxTxSize uint32,
	blockTime time.Duration,
	pendingID types.Proposal,
	pendingTime time.Duration,
	pendedCh chan *Pended,
) *PendingOpts {
	return &PendingOpts{
		maxTxSize:   maxTxSize,
		blockTime:   blockTime,
		pendingID:   pendingID,
		pendingTime: pendingTime,
		pendedCh:    pendedCh,
	}
}
