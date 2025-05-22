package mempool

import (
	"log"
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
	assert.Nil(t, err)

	pn := p.pendings["Pending"]
	assert.Equal(t, types.Topic("Pending"), pn.pendingID)
	assert.Equal(t, time.Second, pn.blockTime)
	assert.Equal(t, uint32(50000), pn.maxTransactionSize)
	assert.Equal(t, 3*time.Second, pn.pendingTime)
	assert.Equal(t, 0, len(pn.scheduledBlockHeight))

	tx1Hash := util.RandomHash()
	tx1 := transaction.NewTransaction(tx1Hash, "P")

	assert.Nil(t, pn.PushTx(tx1))

	log.Println(pn.transactions[tx1.GetHashString()])

	tx2_Hash := util.RandomHash()
	tx2 := transaction.NewTransaction(tx2_Hash, "P")

	assert.Nil(t, pn.PushTx(tx2))

	log.Println(pn.transactions[tx2.GetHashString()])

	time.Sleep(time.Second)

	assert.Equal(t, 2, pn.Len())
}
