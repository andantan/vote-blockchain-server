package gRPC

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

func (vs *VoteSubmit) Fragmentation() (types.Proposal, *transaction.Transaction) {
	plainText := fmt.Sprintf("\"%s\"|\"%s\"|\"%s\"", vs.UserHash.String(), vs.Topic, vs.Option)
	digest := sha256.Sum256([]byte(plainText))

	return vs.Topic, transaction.NewTransaction(types.Hash(digest), vs.Option, time.Now().UnixNano())
}
