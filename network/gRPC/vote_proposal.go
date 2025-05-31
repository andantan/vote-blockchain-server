package gRPC

import (
	"time"

	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_proposal_message"
	"github.com/andantan/vote-blockchain-server/types"
)

// Mapping request - response
type VoteProposal struct {
	Topic      types.Topic
	Duration   time.Duration
	ResponseCh chan *VoteProposalResponse
}

func NewVoteProposal(t *vote_proposal_message.VoteProposalRequest) *VoteProposal {
	return &VoteProposal{
		Topic:    types.Topic(t.GetTopic()),
		Duration: time.Duration(t.GetDuration()) * time.Minute,
	}
}

type VoteProposalResponse struct {
	Status  string
	Message string
	Success bool
}

func NewVoteProposalResponse(status, message string, success bool) *VoteProposalResponse {
	return &VoteProposalResponse{
		Status:  status,
		Message: message,
		Success: success,
	}
}

func GetSuccessVoteProposal(message string) *VoteProposalResponse {
	return NewVoteProposalResponse("SUCCESS", message, true)
}

func GetErrorVoteProposal(message string) *VoteProposalResponse {
	return NewVoteProposalResponse("ERROR", message, false)
}

func (p *VoteProposalResponse) GetTopicResponse() *vote_proposal_message.VoteProposalResponse {
	return &vote_proposal_message.VoteProposalResponse{
		Status:  p.Status,
		Message: p.Message,
		Success: p.Success,
	}
}
