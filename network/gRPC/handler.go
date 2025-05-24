package gRPC

import (
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

func (v *PreTxVote) Fragmentation() (types.Topic, *transaction.Transaction) {
	tx := transaction.NewTransaction(v.Hash, v.Option, time.Now().UnixNano())

	return v.Topic, tx
}
