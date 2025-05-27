package transaction

import (
	"fmt"

	"github.com/andantan/vote-blockchain-server/types"
)

type Transaction struct {
	hash      types.Hash
	option    string
	timeStamp int64
}

func NewTransaction(hash types.Hash, option string, timeStamp int64) *Transaction {
	return &Transaction{
		hash:      hash,
		option:    option,
		timeStamp: timeStamp,
	}
}

func (tx *Transaction) GetHash() types.Hash {
	return tx.hash
}

func (tx *Transaction) GetHashString() string {
	return tx.hash.String()
}

func (tx *Transaction) GetOption() string {
	return tx.option
}

func (tx *Transaction) GetTimeStamp() int64 {
	return tx.timeStamp
}

// Return "Hash|Option"
func (tx *Transaction) Serialize() string {
	s := fmt.Sprintf("%s|%s", tx.hash.String(), tx.option)

	return s
}
