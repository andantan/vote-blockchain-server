package mempool

import (
	"fmt"
	"log"
	"time"

	"github.com/andantan/vote-blockchain-server/types"
)

type MemPool struct {
	BlockTime time.Duration
	MaxTxSize uint32
	pendings  map[types.Topic]*Pending
}

func NewMemPool(blockTime time.Duration, maxTxSize uint32) *MemPool {
	return &MemPool{
		BlockTime: blockTime,
		MaxTxSize: maxTxSize,
		pendings:  make(map[types.Topic]*Pending),
	}
}

func (p *MemPool) AddPending(pendingId types.Topic, pendingTime time.Duration) error {
	if p.IsOpen(pendingId) {
		return fmt.Errorf("topic(%s) already opened pending", pendingId)
	}

	open := NewPending(pendingTime, p.BlockTime, p.MaxTxSize, pendingId)

	p.AllocatePending(pendingId, open)

	log.Printf("topic(%s) pending success, duration (%s)", open.pendingID, open.pendingTime)

	return nil
}

// Check Pending is opened
func (p *MemPool) IsOpen(pendingId types.Topic) bool {
	_, ok := p.pendings[pendingId]

	return ok
}

func (p *MemPool) AllocatePending(pendingId types.Topic, open *Pending) {
	p.pendings[pendingId] = open
}
