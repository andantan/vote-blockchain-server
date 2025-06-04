package gRPC

import (
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

func (vs *VoteSubmit) Fragmentation() (types.Proposal, *transaction.Transaction) {
	return vs.Topic, transaction.NewTransaction(vs.Hash, vs.Option, time.Now().UnixNano())
}
