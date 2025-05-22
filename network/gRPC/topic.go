package gRPC

import (
	"time"

	"github.com/andantan/vote-blockchain-server/network/gRPC/topic_message"
	"github.com/andantan/vote-blockchain-server/types"
)

// Mapping request - response
type PreTxTopic struct {
	Topic      types.Topic
	Duration   time.Duration
	ResponseCh chan *PostTxTopic
}

func GetPreTxTopic(t *topic_message.TopicRequest) *PreTxTopic {
	return &PreTxTopic{
		Topic:    types.Topic(t.GetTopic()),
		Duration: time.Duration(t.GetDuration()) * time.Minute,
	}
}

type PostTxTopic struct {
	Status  string
	Message string
	Success bool
}

func GetPostTxTopic(status, message string, success bool) *PostTxTopic {
	return &PostTxTopic{
		Status:  status,
		Message: message,
		Success: success,
	}
}

func (p *PostTxTopic) GetTopicResponse() *topic_message.TopicResponse {
	return &topic_message.TopicResponse{
		Status:  p.Status,
		Message: p.Message,
		Success: p.Success,
	}
}
