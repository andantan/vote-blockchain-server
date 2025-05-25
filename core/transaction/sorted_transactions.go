package transaction

import (
	"sort"

	"github.com/andantan/vote-blockchain-server/types"
)

type SortedTxx struct {
	txx []*Transaction
}

func NewSortedTxx(txMap map[string]*Transaction) *SortedTxx {

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

func (s *SortedTxx) GetHashSlice() []types.Hash {
	hashes := make([]types.Hash, 0, len(s.GetTxx()))

	for _, tx := range s.GetTxx() {
		hashes = append(hashes, tx.GetHash())
	}

	return hashes
}
