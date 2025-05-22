package mempool

import (
	"testing"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
	"github.com/stretchr/testify/assert"
)

func TestMempool(t *testing.T) {
	p := NewMemPool(5*time.Second, uint32(50000))

	assert.False(t, p.IsOpen("Pending1"))

	err := p.AddPending("Pending1", 5*time.Hour)

	assert.Nil(t, err)
	assert.True(t, p.IsOpen("Pending1"))
	assert.Equal(t, 5*time.Second, p.BlockTime)
	assert.Equal(t, uint32(50000), p.MaxTxSize)
}

func TestPending(t *testing.T) {
	p := NewMemPool(time.Second, uint32(50000))

	err := p.AddPending("Pending", 3*time.Second)
	pn := p.pendings["Pending"]

	assert.Nil(t, err)
	assert.Equal(t, types.Topic("Pending"), pn.pendingID)
	assert.Equal(t, time.Second, pn.blockTime)
	assert.Equal(t, uint32(50000), pn.maxTransactionSize)
	assert.Equal(t, 3*time.Second, pn.pendingTime)
	assert.Equal(t, 0, len(pn.scheduledBlockHeight))

	tx := transaction.Transaction{
		Hash:   util.RandomHash(),
		Option: "P",
	}

	assert.Nil(t, pn.CommitTx(&tx))

	tx2 := transaction.Transaction{
		Hash:   util.RandomHash(),
		Option: "P",
	}

	assert.Nil(t, pn.CommitTx(&tx2))

}
