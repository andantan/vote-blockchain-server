package gRPC

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

// Return (*core.Transaction, electionId) from Vote
func (v *Vote) Fragmentation() (*transaction.Transaction, types.VotingID) {
	return &transaction.Transaction{
		VoteHash:   v.VoteHash,
		VoteOption: v.VoteOption,
	}, v.VoteId
}
