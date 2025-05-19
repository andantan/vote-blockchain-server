package vote

import (
	"testing"

	"github.com/andantan/vote-blockchain-server/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestXxx(t *testing.T) {
	rv := util.RandomVote()

	req := VoteRequest{
		VoteHash:   rv.VoteHash.String(),
		VoteOption: rv.VoteOption,
		ElectionId: rv.ElectionId,
	}

	data, err := proto.Marshal(&req)

	assert.Nil(t, err)

	ureq := new(VoteRequest)

	proto.Unmarshal(data, ureq)

	assert.Equal(t, req.VoteHash, ureq.VoteHash)
	assert.Equal(t, req.VoteOption, ureq.VoteOption)
	assert.Equal(t, req.ElectionId, ureq.ElectionId)
}
