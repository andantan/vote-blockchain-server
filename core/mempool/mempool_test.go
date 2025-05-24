package mempool

import (
	"testing"
	"time"

	"github.com/andantan/vote-blockchain-server/core/signal"
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
	pendingName := types.Topic("pending")

	mp := NewMemPool(5*time.Second, uint32(50000))
	mp.SetChannel(nil)

	err := mp.AddPending(pendingName, 10*time.Second)
	assert.Nil(t, err)

	pn := mp.pendings[pendingName]

	assert.True(t, mp.IsOpen(pendingName))
	assert.Equal(t, 0, pn.Len())

	tx1Hash := util.RandomHash()
	tx1 := randomTx(tx1Hash, "P")

	tx2_Hash := util.RandomHash()
	tx2 := randomTx(tx2_Hash, "P")

	assert.Nil(t, pn.PushTx(tx1))
	assert.Nil(t, mp.CommitTransaction(pendingName, tx2))

	time.Sleep(time.Second)

	assert.Equal(t, 2, pn.Len())
	atx1 := pn.seekTx(tx1.GetHashString())
	atx2 := pn.seekTx(tx2.GetHashString())

	assert.NotNil(t, atx1)
	assert.NotNil(t, atx2)

	t.Log(atx1)
	t.Log(atx2)

	assert.Equal(t, tx1Hash, atx1.GetHash())
	assert.Equal(t, tx2_Hash, atx2.GetHash())
	assert.Equal(t, atx1.GetOption(), atx2.GetOption())

	sync := signal.NewPendingClosing(pendingName, pn.wg, 300*time.Millisecond)

	startTime := time.Now()

	sync.Add(1)
	mp.closeCh <- sync
	sync.Wait()
	elapsedTime := time.Since(startTime)

	t.Logf("%s", elapsedTime)
	assert.False(t, mp.IsOpen(pendingName))
}

func randomTx(hash types.Hash, option string) *transaction.Transaction {
	return transaction.NewTransaction(hash, option, time.Now().UnixNano())
}
