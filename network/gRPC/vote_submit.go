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
	UserHash   types.Hash
	Option     string
	Topic      types.Proposal
	ResponseCh chan *VoteSubmitResponse
}

func NewVoteSubmit(v *vote_submit_message.SubmitBallotTransactionRequest) (*VoteSubmit, error) {
	s := v.GetUserHash()

	if len(s) != HASH_STRING_SIZE {
		msg := fmt.Sprintf("Given string hash with length %d should be %d", len(s), HASH_STRING_SIZE)

		return nil, werror.NewWrappedError("INVALID_HASH_LENGTH", msg, nil)
	}

	h, err := types.HashFromHashString(s)

	if err != nil {
		return nil, werror.NewWrappedError("DECODE_ERROR", err.Error(), err)
	}

	return &VoteSubmit{
		UserHash: h,
		Option:   v.GetOption(),
		Topic:    types.Proposal(v.GetTopic()),
	}, nil
}

type VoteSubmitResponse struct {
	VoteHash string
	Success  bool
	Status   string
}

func NewVoteSubmitResponse(status, voteHash string, success bool) *VoteSubmitResponse {
	return &VoteSubmitResponse{
		VoteHash: voteHash,
		Success:  success,
		Status:   status,
	}
}

func NewSuccessVoteSubmitResponse(hash types.Hash) *VoteSubmitResponse {
	return NewVoteSubmitResponse("SUBMITTED", hash.String(), true)
}

func NewErrorVoteSubmitResponse(err error) *VoteSubmitResponse {
	if err == nil {
		return NewVoteSubmitResponse("INTERNAL_ERROR", types.NilHash().String(), false)
	}

	if werr, ok := err.(*werror.WrappedError); ok {
		return NewVoteSubmitResponse(werr.Code, types.NilHash().String(), false)
	}

	return NewVoteSubmitResponse("UNKNOWN_ERROR", types.NilHash().String(), false)
}

// TODO change proto message type name
func (p *VoteSubmitResponse) GetVoteResponse() *vote_submit_message.SubmitBallotTransactionResponse {
	return &vote_submit_message.SubmitBallotTransactionResponse{
		VoteHash: p.VoteHash,
		Success:  p.Success,
		Status:   p.Status,
	}
}
