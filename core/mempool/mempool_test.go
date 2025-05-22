package mempool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMempool(t *testing.T) {
	p := NewMemPool(5*time.Second, uint32(50000))

	assert.False(t, p.IsOpen("Pending1"))
	assert.Nil(t, p.AddPending("Pending1", 5*time.Hour))
	assert.True(t, p.IsOpen("Pending1"))
	assert.Equal(t, 5*time.Second, p.BlockTime)
	assert.Equal(t, uint32(50000), p.MaxTxSize)
}
