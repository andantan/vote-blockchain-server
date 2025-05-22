package gRPC

import (
	"fmt"

	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_message"
	"github.com/andantan/vote-blockchain-server/types"
)

type PreTxVote struct {
	Hash   types.Hash
	Option string
	Topic  types.Topic
}

func GetPreTxVote(v *vote_message.VoteRequest) (*PreTxVote, error) {
	s := v.GetHash()

	if len(s) != 64 {
		return nil, fmt.Errorf("given string hash with length %d should be 64", len(s))
	}

	h, err := types.HashFromHashString(s)

	if err != nil {
		// return nil, fmt.Errorf("given string voteHash (%s) does not satisfy the hash string", s)
		return nil, fmt.Errorf("%s", err.Error())
	}

	return &PreTxVote{
		Hash:   h,
		Option: v.GetOption(),
		Topic:  types.Topic(v.GetTopic()),
	}, nil
}

type PostTxVote struct {
	Status  string
	Message string
	Success bool
}

func GetPostTxVote(status, message string, success bool) *PostTxVote {
	return &PostTxVote{
		Status:  status,
		Message: message,
		Success: success,
	}
}

func (p *PostTxVote) GetVoteResponse() *vote_message.VoteResponse {
	return &vote_message.VoteResponse{
		Status:  p.Status,
		Message: p.Message,
		Success: p.Success,
	}
}
