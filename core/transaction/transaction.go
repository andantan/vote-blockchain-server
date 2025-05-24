package transaction

import (
	"fmt"
	"sort"

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

func (tx *Transaction) GetTimeStamp() int64 {
	return tx.timeStamp
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

// Return "Hash|Option|timestamp"
func (tx *Transaction) Serialize() string {
	s := fmt.Sprintf("%s|%s|%d", tx.hash.String(), tx.option, tx.timeStamp)

	return s
}

type SortedTxx struct {
	txx []*Transaction
}

func NewTxMapSorter(txMap map[string]*Transaction) *SortedTxx {
	txx := make([]*Transaction, len(txMap))

	i := 0

	for _, val := range txMap {
		txx[i] = val
		i++
	}

	s := &SortedTxx{
		txx: txx,
	}

	sort.Sort(s)

	return s
}

func (s *SortedTxx) Len() int {
	return len(s.txx)
}

func (s *SortedTxx) Swap(i, j int) {
	s.txx[i], s.txx[j] = s.txx[j], s.txx[i]
}

func (s *SortedTxx) Less(i, j int) bool {
	// Alphabetical order
	if s.txx[i].timeStamp == s.txx[j].timeStamp {
		return s.txx[i].GetHashString() < s.txx[j].GetHashString()
	}
	return s.txx[i].timeStamp < s.txx[j].timeStamp
}

func (s *SortedTxx) GetTxx() []*Transaction {
	return s.txx
}
