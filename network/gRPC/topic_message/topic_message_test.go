package topic_message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestVoteMessage(t *testing.T) {

	req := TopicRequest{
		Topic:    "2025-보건의료여론조사",
		Duration: 7200,
	}

	data, err := proto.Marshal(&req)

	assert.Nil(t, err)

	ureq := new(TopicRequest)

	proto.Unmarshal(data, ureq)

	assert.Equal(t, req.Topic, ureq.Topic)
	assert.Equal(t, req.Duration, ureq.Duration)
}
