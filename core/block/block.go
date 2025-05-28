package block

import (
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type Block struct {
	*Header      `json:"header"`
	BlockHash    types.Hash                 `json:"block_hash"`
	Transactions []*transaction.Transaction `json:"transactions"`
}

func NewBlock(h *Header, txx []*transaction.Transaction) *Block {
	return &Block{
		Header:       h,
		BlockHash:    h.Hash(),
		Transactions: txx,
	}
}

func NewBlockFromPrevHeader(prevHeader *Header, pb *ProtoBlock) *Block {
	header := &Header{
		VotingID:      pb.VotingID,
		MerkleRoot:    pb.MerkleRoot,
		Height:        prevHeader.Height + 1,
		PrevBlockHash: prevHeader.Hash(),
	}

	return NewBlock(header, pb.txx)
}

func GenesisBlock() *Block {

	gpb := genesisProtoBlock()

	gh := &Header{
		VotingID:      "GENESIS",
		MerkleRoot:    types.FFHashCompact(),
		Height:        0,
		PrevBlockHash: types.ZeroHashCompact(),
	}

	return NewBlock(gh, gpb.txx)
}
