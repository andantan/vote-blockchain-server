package block

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type ProtoBlock struct {
	VotingID   types.Topic
	MerkleRoot types.Hash
	txx        []*transaction.Transaction
}

// Preprocessing for create block
func NewProtoBlock(ID types.Topic, txMap map[string]*transaction.Transaction) *ProtoBlock {
	stx := transaction.NewSortedTxx(txMap)

	return &ProtoBlock{
		VotingID:   ID,
		MerkleRoot: CalculateMerkleRoot(stx),
		txx:        stx.GetTxx(),
	}
}

func genesisProtoBlock() *ProtoBlock {
	return &ProtoBlock{
		VotingID:   "GENESIS",
		MerkleRoot: types.NilHash(),
		txx:        []*transaction.Transaction{},
	}
}
