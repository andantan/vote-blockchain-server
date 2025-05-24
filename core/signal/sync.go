package signal

import (
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/types"
)

type PendingClosing struct {
	*sync.WaitGroup
	topic    types.Topic
	syncTime time.Duration
}

func NewPendingClosing(topic types.Topic, wg *sync.WaitGroup, t time.Duration) *PendingClosing {
	return &PendingClosing{
		WaitGroup: wg,
		topic:     topic,
		syncTime:  t,
	}
}

func (c *PendingClosing) GetTopic() types.Topic {
	return c.topic
}

func (c *PendingClosing) GetSyncTime() time.Duration {
	return c.syncTime
}
