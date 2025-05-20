package mempool

import (
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type Pending struct {
	transactions         []*transaction.Transaction // Votes [ Transaction ]
	pendingTime          time.Duration              // Vote duration
	blockTime            time.Duration              // Block Time (system)
	maxTransactionSize   uint32                     // Tx size (system)
	scheduledBlockHeight []uint64                   // Pended block heights
	pendingId            types.ElectionID           // ElectionID
}

func NewPending(pendingTime, blockTime time.Duration,
	maxTxSize uint32, pendingId string) *Pending {
	return &Pending{
		transactions:         []*transaction.Transaction{},
		pendingTime:          pendingTime,
		blockTime:            blockTime,
		maxTransactionSize:   maxTxSize,
		scheduledBlockHeight: []uint64{},
		pendingId:            types.ElectionID(pendingId),
	}
}
