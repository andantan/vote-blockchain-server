package vote_message

import (
	"testing"

	"github.com/andantan/vote-blockchain-server/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestVoteMessage(t *testing.T) {
	rv := util.RandomVote()

	req := VoteRequest{
		VoteHash:   rv.VoteHash.String(),
		VoteOption: rv.VoteOption,
		VoteId:     string(rv.VoteId),
	}

	data, err := proto.Marshal(&req)

	assert.Nil(t, err)

	ureq := new(VoteRequest)

	proto.Unmarshal(data, ureq)

	assert.Equal(t, req.VoteHash, ureq.VoteHash)
	assert.Equal(t, req.VoteOption, ureq.VoteOption)
	assert.Equal(t, req.VoteId, ureq.VoteId)
}
