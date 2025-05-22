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
		Hash:   rv.VoteHash.String(),
		Option: rv.VoteOption,
		Topic:  string(rv.VoteId),
	}

	data, err := proto.Marshal(&req)

	assert.Nil(t, err)

	ureq := new(VoteRequest)

	proto.Unmarshal(data, ureq)

	assert.Equal(t, req.Hash, ureq.Hash)
	assert.Equal(t, req.Option, ureq.Option)
	assert.Equal(t, req.Topic, ureq.Topic)
}
