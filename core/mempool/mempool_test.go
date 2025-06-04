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
	t.Log(util.GreenString("+-+-+-+-+-+-+-+-+- mempool_test.go::TestPending \"RUN\" -+-+-+-+-+-+-+-+-+"))

	pendingName := types.Proposal("pending")

	mp := NewMemPool(3*time.Second, uint32(50000))

	err := mp.AddPending(pendingName, 5*time.Second)
	assert.Nil(t, err)

	pn := mp.pendings[pendingName]

	assert.True(t, mp.IsOpen(pendingName))
	assert.Equal(t, 0, pn.Len())

	tx1Hash := util.RandomHash()
	tx1 := randomTx(tx1Hash, "P")

	tx2_Hash := util.RandomHash()
	tx2 := randomTx(tx2_Hash, "P")

	t.Logf(util.CyanString("Tx1.Hash: %v"), tx1.GetHash())
	t.Logf(util.CyanString("Tx1.Option: %s"), tx1.GetOption())
	t.Logf(util.CyanString("Tx1.TimeStamp: %d"), tx1.GetTimeStamp())
	t.Logf(util.CyanString("Tx2.Hash: %v"), tx2.GetHash())
	t.Logf(util.CyanString("Tx2.Option: %s"), tx2.GetOption())
	t.Logf(util.CyanString("Tx2.TimeStamp: %d"), tx2.GetTimeStamp())

	assert.Nil(t, pn.PushTx(tx1))
	assert.Nil(t, mp.CommitTransaction(pendingName, tx2))

	time.Sleep(time.Second)

	assert.Equal(t, 2, pn.Len())
	atx1 := pn.seekTx(tx1.GetHashString())
	atx2 := pn.seekTx(tx2.GetHashString())

	assert.NotNil(t, atx1)
	assert.NotNil(t, atx2)

	t.Logf(util.MagentaString("aTx1.Hash: %v"), atx1.GetHash())
	t.Logf(util.MagentaString("aTx1.Option: %s"), atx1.GetOption())
	t.Logf(util.MagentaString("aTx1.TimeStamp: %d"), atx1.GetTimeStamp())
	t.Logf(util.MagentaString("aTx2.Hash: %v"), atx2.GetHash())
	t.Logf(util.MagentaString("aTx2.Option: %s"), atx2.GetOption())
	t.Logf(util.MagentaString("aTx2.TimeStamp: %d"), atx2.GetTimeStamp())

	assert.Equal(t, tx1Hash, atx1.GetHash())
	assert.Equal(t, tx2_Hash, atx2.GetHash())
	assert.Equal(t, atx1.GetOption(), atx2.GetOption())
	t.Log(util.GreenString("+-+-+-+-+-+-+-+-+- mempool_test.go::TestPending \"END\" -+-+-+-+-+-+-+-+-+"))
}

func randomTx(hash types.Hash, option string) *transaction.Transaction {
	return transaction.NewTransaction(hash, option, time.Now().UnixNano())
}
