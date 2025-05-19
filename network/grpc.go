package network

import (
	"log"

	"github.com/andantan/vote-blockchain-server/core"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/vote"
)

type Vote struct {
	VoteHash   types.Hash
	VoteOption string
	ElectionId string
}

func GetVoteFromgRPCRequest(v *vote.VoteRequest) Vote {
	s := v.GetVoteHash()
	h, err := types.HashFromHashString(s)

	if err != nil {
		log.Fatalf("given string voteHash (%s) does not satisfy the hash string.", s)

		return Vote{}
	}

	return Vote{
		VoteHash:   h,
		VoteOption: v.GetVoteOption(),
		ElectionId: v.GetElectionId(),
	}
}

func (v *Vote) Fragmentation() (*core.Transaction, string) {
	return &core.Transaction{
		VoteHash:   v.VoteHash,
		VoteOption: v.VoteOption,
	}, v.ElectionId
}
