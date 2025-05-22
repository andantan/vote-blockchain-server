package gRPC

import (
	"log"

	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_message"
	"github.com/andantan/vote-blockchain-server/types"
)

type PreTxVote struct {
	Hash   types.Hash
	Option string
	Topic  types.Topic
}

func GetVoteFromVoteMessage(v *vote_message.VoteRequest) *PreTxVote {
	s := v.GetHash()
	h, err := types.HashFromHashString(s)

	if err != nil {
		log.Fatalf("given string voteHash (%s) does not satisfy the hash string.", s)

		return &PreTxVote{}
	}

	return &PreTxVote{
		Hash:   h,
		Option: v.GetOption(),
		Topic:  types.Topic(v.GetTopic()),
	}
}

type PostTxVote struct {
	Status  string
	message string
	success bool
}
