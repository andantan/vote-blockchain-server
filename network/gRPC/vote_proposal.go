package gRPC

import (
	"fmt"
	"time"

	werror "github.com/andantan/vote-blockchain-server/error"
	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_proposal_message"
	"github.com/andantan/vote-blockchain-server/types"
)

// Mapping request - response
type VoteProposal struct {
	Proposal   types.Proposal
	Duration   time.Duration
	ResponseCh chan *VoteProposalResponse
}

func NewVoteProposal(p *vote_proposal_message.OpenProposalRequest) *VoteProposal {
	return &VoteProposal{
		Proposal: types.Proposal(p.GetTopic()),
		Duration: time.Duration(p.GetDuration()) * time.Minute,
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

func NewSuccessVoteProposalResponse(proposal types.Proposal, duration time.Duration) *VoteProposalResponse {
	msg := fmt.Sprintf("Proposal '%s' is now open for pending submissions. Duration: %s.", proposal, duration)

	return NewVoteProposalResponse("OK", msg, true)
}

func NewErrorVoteProposalResponse(err error) *VoteProposalResponse {
	if err == nil {
		return NewVoteProposalResponse("INTERNAL_ERROR", "Unexpected error occurred (nil error provided).", false)
	}

	if werr, ok := err.(*werror.WrappedError); ok {
		return NewVoteProposalResponse(werr.Code, werr.Message, false)
	}

	return NewVoteProposalResponse("UNKNOWN_ERROR", fmt.Sprintf("Unexpected error occurred: %v", err), false)
}

func (p *VoteProposalResponse) GetTopicResponse() *vote_proposal_message.OpenProposalResponse {
	return &vote_proposal_message.OpenProposalResponse{
		Status:  p.Status,
		Message: p.Message,
		Success: p.Success,
	}
}
