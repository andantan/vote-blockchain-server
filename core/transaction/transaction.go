package transaction

import (
	"fmt"

	"github.com/andantan/vote-blockchain-server/types"
)

type Transaction struct {
	Hash      types.Hash `json:"hash"`
	Option    string     `json:"option"`
	TimeStamp int64      `json:"time_stamp"`
}

func NewTransaction(voteHash types.Hash, option string, timeStamp int64) *Transaction {
	return &Transaction{
		Hash:      voteHash,
		Option:    option,
		TimeStamp: timeStamp,
	}
}

func (tx *Transaction) GetHash() types.Hash {
	return tx.Hash
}

func (tx *Transaction) GetHashString() string {
	return tx.Hash.String()
}

func (tx *Transaction) GetOption() string {
	return tx.Option
}

func (tx *Transaction) GetTimeStamp() int64 {
	return tx.TimeStamp
}

// Return "Hash|Option"
func (tx *Transaction) Serialize() string {
	s := fmt.Sprintf("%s|%s", tx.Hash.String(), tx.Option)

	return s
}
