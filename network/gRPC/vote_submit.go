package gRPC

import (
	"fmt"

	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_message"
	"github.com/andantan/vote-blockchain-server/types"
)

const (
	HASH_STRING_SIZE = 64
)

type VoteSubmit struct {
	Hash       types.Hash
	Option     string
	Topic      types.Topic
	ResponseCh chan *VoteSubmitResponse
}

func NewVoteSubmit(v *vote_message.VoteRequest) (*VoteSubmit, error) {
	s := v.GetHash()

	if len(s) != HASH_STRING_SIZE {
		return nil, fmt.Errorf("given string hash with length %d should be 64", len(s))
	}

	h, err := types.HashFromHashString(s)

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return &VoteSubmit{
		Hash:   h,
		Option: v.GetOption(),
		Topic:  types.Topic(v.GetTopic()),
	}, nil
}

type VoteSubmitResponse struct {
	Status  string
	Message string
	Success bool
}

func NewVoteSubmitResponse(status, message string, success bool) *VoteSubmitResponse {
	return &VoteSubmitResponse{
		Status:  status,
		Message: message,
		Success: success,
	}
}

func GetSuccessSubmitVote(message string) *VoteSubmitResponse {
	return NewVoteSubmitResponse("SUCCESS", message, true)
}

func GetErrorSubmitVote(message string) *VoteSubmitResponse {
	return NewVoteSubmitResponse("ERROR", message, false)
}

func (p *VoteSubmitResponse) GetVoteResponse() *vote_message.VoteResponse {
	return &vote_message.VoteResponse{
		Status:  p.Status,
		Message: p.Message,
		Success: p.Success,
	}
}
