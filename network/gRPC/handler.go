package gRPC

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

// Return (*core.Transaction, electionId) from Vote
func (v *PreTxVote) Fragmentation() (types.Topic, *transaction.Transaction) {
	tx := transaction.NewTransaction(v.Hash, v.Option)

	return v.Topic, tx
}
