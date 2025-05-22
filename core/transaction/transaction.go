package transaction

import (
	"fmt"

	"github.com/andantan/vote-blockchain-server/types"
)

type Transaction struct {
	hash   types.Hash
	option string
}

func NewTransaction(hash types.Hash, option string) *Transaction {
	return &Transaction{
		hash:   hash,
		option: option,
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

// Return "Hash|Option"
func (tx *Transaction) Serialize() string {
	s := fmt.Sprintf("%s|%s", tx.hash.String(), tx.option)

	return s
}
