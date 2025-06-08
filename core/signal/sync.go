package signal

import (
	"sync"

	"github.com/andantan/vote-blockchain-server/types"
)

type PendingClosing struct {
	wg    *sync.WaitGroup
	topic types.Proposal
}

func NewPendingClosing(topic types.Proposal) *PendingClosing {
	return &PendingClosing{
		wg:    &sync.WaitGroup{},
		topic: topic,
	}
}

func (c *PendingClosing) GetTopic() types.Proposal {
	return c.topic
}

func (c *PendingClosing) Add(delta int) {
	c.wg.Add(delta)
}

func (c *PendingClosing) Done() {
	c.wg.Done()
}

func (c *PendingClosing) Wait() {
	c.wg.Wait()
}
