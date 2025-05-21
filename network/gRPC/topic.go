package gRPC

import (
	"time"

	"github.com/andantan/vote-blockchain-server/network/gRPC/topic_message"
	"github.com/andantan/vote-blockchain-server/types"
)

type Topic struct {
	TopicId       types.TopicID
	TopicDuration time.Duration
}

func GetTopicFromTopicMessage(t *topic_message.TopicRequest) Topic {
	return Topic{
		TopicId:       types.TopicID(t.GetTopicId()),
		TopicDuration: time.Duration(t.GetTopicDuration()) * time.Minute,
	}
}
