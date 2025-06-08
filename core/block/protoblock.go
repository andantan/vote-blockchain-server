package block

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type ProtoBlock struct {
	VotingID   types.Proposal
	MerkleRoot types.Hash
	txx        []*transaction.Transaction
}

// Preprocessing for create block
func NewProtoBlock(proposal types.Proposal, txMap map[string]*transaction.Transaction) *ProtoBlock {
	stx := transaction.NewSortedTxx(txMap)

	return &ProtoBlock{
		VotingID:   proposal,
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

func (pb *ProtoBlock) Len() int {
	return len(pb.txx)
}
