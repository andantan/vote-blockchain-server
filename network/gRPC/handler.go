package gRPC

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

// TODO: Arbitary User's Salt for Hashing
func (vs *VoteSubmit) Fragmentation() (types.Proposal, *transaction.Transaction) {
	plainText := fmt.Sprintf("\"%s\"|\"%s\"|\"%s\"|\"%s\"", vs.UserHash.String(), vs.Topic, vs.Option, vs.Salt)
	digest := sha256.Sum256([]byte(plainText))

	return vs.Topic, transaction.NewTransaction(types.Hash(digest), vs.Option, time.Now().UnixNano())
}
