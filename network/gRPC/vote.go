package gRPC

import (
	"log"

	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_message"
	"github.com/andantan/vote-blockchain-server/types"
)

type Vote struct {
	VoteHash   types.Hash
	VoteOption string
	VoteId     types.VotingID
}

func GetVoteFromVoteMessage(v *vote_message.VoteRequest) Vote {
	s := v.GetVoteHash()
	h, err := types.HashFromHashString(s)

	if err != nil {
		log.Fatalf("given string voteHash (%s) does not satisfy the hash string.", s)

		return Vote{}
	}

	return Vote{
		VoteHash:   h,
		VoteOption: v.GetVoteOption(),
		VoteId:     types.VotingID(v.GetVoteId()),
	}
}
