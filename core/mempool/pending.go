package mempool

import (
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type Pending struct {
	transactions         []*transaction.Transaction // Votes [ Transaction ]
	transactionCH        chan *transaction.Transaction
	pendingTime          time.Duration // Vote duration
	blockTime            time.Duration // Block Time (system)
	maxTransactionSize   uint32        // Tx size (system)
	scheduledBlockHeight []uint64      // Pended block heights
	pendingID            types.Topic   // TopicID
}

func NewPending(pendingTime, blockTime time.Duration,
	maxTxSize uint32, pendingId types.Topic) *Pending {
	return &Pending{
		transactions:         []*transaction.Transaction{},
		transactionCH:        make(chan *transaction.Transaction),
		pendingTime:          pendingTime,
		blockTime:            blockTime,
		maxTransactionSize:   maxTxSize,
		scheduledBlockHeight: []uint64{},
		pendingID:            pendingId,
	}
}
