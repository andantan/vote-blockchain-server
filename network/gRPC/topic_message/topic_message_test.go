package topic_message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestVoteMessage(t *testing.T) {

	req := TopicRequest{
		TopicId:       "2025-보건의료여론조사",
		TopicDuration: 7200,
	}

	data, err := proto.Marshal(&req)

	assert.Nil(t, err)

	ureq := new(TopicRequest)

	proto.Unmarshal(data, ureq)

	assert.Equal(t, req.TopicId, ureq.TopicId)
	assert.Equal(t, req.TopicDuration, ureq.TopicDuration)
}
