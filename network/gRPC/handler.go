package gRPC

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

// Return (*core.Transaction, electionId) from Vote
func (v *PreTxVote) Fragmentation() (*transaction.Transaction, types.Topic) {
	return &transaction.Transaction{
		Hash:   v.Hash,
		Option: v.Option,
	}, v.Topic
}
