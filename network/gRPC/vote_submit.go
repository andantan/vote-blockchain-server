package gRPC

import (
	"fmt"

	werror "github.com/andantan/vote-blockchain-server/error"
	"github.com/andantan/vote-blockchain-server/network/gRPC/vote_submit_message"
	"github.com/andantan/vote-blockchain-server/types"
)

const (
	HASH_STRING_SIZE = 64
)

type VoteSubmit struct {
	Hash       types.Hash
	Option     string
	Topic      types.Proposal
	ResponseCh chan *VoteSubmitResponse
}

func NewVoteSubmit(v *vote_submit_message.VoteSubmitRequest) (*VoteSubmit, error) {
	s := v.GetHash()

	if len(s) != HASH_STRING_SIZE {
		msg := fmt.Sprintf("Given string hash with length %d should be %d", len(s), HASH_STRING_SIZE)

		return nil, werror.NewWrappedError("INVALID_HASH_LENGTH", msg, nil)
	}

	h, err := types.HashFromHashString(s)

	if err != nil {
		return nil, werror.NewWrappedError("DECODE_ERROR", err.Error(), err)
	}

	return &VoteSubmit{
		Hash:   h,
		Option: v.GetOption(),
		Topic:  types.Proposal(v.GetTopic()),
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

func NewSuccessVoteSubmitResponse(proposal types.Proposal, hash types.Hash) *VoteSubmitResponse {
	msg := fmt.Sprintf("Vote for proposal '%s' has been successfully submitted with transaction hash '%s'.", proposal, hash)

	return NewVoteSubmitResponse("SUBMITTED", msg, true)
}

func NewErrorVoteSubmitResponse(err error) *VoteSubmitResponse {
	if err == nil {
		return NewVoteSubmitResponse("INTERNAL_ERROR", "Unexpected error occurred (nil error provided).", false)
	}

	if werr, ok := err.(*werror.WrappedError); ok {
		return NewVoteSubmitResponse(werr.Code, werr.Message, false)
	}

	return NewVoteSubmitResponse("UNKNOWN_ERROR", fmt.Sprintf("Unexpected error occurred: %v", err), false)
}

// TODO change proto message type name
func (p *VoteSubmitResponse) GetVoteResponse() *vote_submit_message.VoteSubmitResponse {
	return &vote_submit_message.VoteSubmitResponse{
		Status:  p.Status,
		Message: p.Message,
		Success: p.Success,
	}
}
