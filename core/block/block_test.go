package block

import (
	"testing"
	"time"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
	"github.com/stretchr/testify/assert"
)

func TestHeaderBytes(t *testing.T) {
	t.Log(util.GreenString("+-+-+-+-+-+-+-+-+- block_test.go::TestHeaderBytes \"RUN\" -+-+-+-+-+-+-+-+-+"))

	randomMerkleRootHash := util.RandomHash()
	randomPrevBlockHash := util.RandomHash()

	t.Logf(util.CyanString("randomMerkleRootHash %v"), randomMerkleRootHash)
	t.Logf(util.CyanString("randomPrevBlockHash %v"), randomPrevBlockHash)

	h1 := &Header{
		VotingID:      "aaa",
		MerkleRoot:    randomMerkleRootHash,
		PrevBlockHash: randomPrevBlockHash,
		Height:        1,
	}

	h2 := &Header{
		VotingID:      "aaa",
		MerkleRoot:    randomMerkleRootHash,
		PrevBlockHash: randomPrevBlockHash,
		Height:        1,
	}

	assert.Equal(t, h1.Bytes(), h2.Bytes())

	h3 := &Header{
		VotingID:      "bbb",
		MerkleRoot:    util.RandomHash(),
		PrevBlockHash: util.RandomHash(),
		Height:        2,
	}

	assert.NotEqual(t, h1.Bytes(), h3.Bytes())
	assert.NotEqual(t, h2.Bytes(), h3.Bytes())

	t.Log(util.GreenString("+-+-+-+-+-+-+-+-+- block_test.go::TestHeaderBytes \"END\" -+-+-+-+-+-+-+-+-+"))
}

func TestPreparedBlock(t *testing.T) {
	t.Log(util.GreenString("+-+-+-+-+-+-+-+-+- block_test.go::TestPreparedBlock \"RUN\" -+-+-+-+-+-+-+-+-+"))

	h1, err1 := types.HashFromHashString("d625bf85840eaf147435f9273d8fb2945b5c2acb84e78e24fee201283c8329a0")
	h2, err2 := types.HashFromHashString("8aefa4c085c4f69914c634f05e2d568ca3711ee291c4159d460094ad67b3523e")
	h3, err3 := types.HashFromHashString("247a776c21fbf4c357426a1477a97d4c421b61b42eced6ccc8202df708de363c")
	h4, err4 := types.HashFromHashString("6112a27532f486ab92d0d9d1689ea7f9fda85bced6185c5bc9dce98bd4967714")
	h5, err5 := types.HashFromHashString("5023e275ebf7d014ec3701a19ca92ccebb97da95a033b502edef2e3364838af3")
	h6, err6 := types.HashFromHashString("8af4c9ce8f6f0ba3efcb536e7a60ef2cf5749f049d0c8c847be5d818a319a72f")

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.Nil(t, err4)
	assert.Nil(t, err5)
	assert.Nil(t, err6)

	time.Sleep(10 * time.Millisecond)
	tx1 := transaction.NewTransaction(h1, "2", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx2 := transaction.NewTransaction(h2, "1", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx3 := transaction.NewTransaction(h3, "3", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx4 := transaction.NewTransaction(h4, "2", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx5 := transaction.NewTransaction(h5, "1", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx6 := transaction.NewTransaction(h6, "3", time.Now().UnixNano())

	txMap := make(map[string]*transaction.Transaction)

	txMap[tx3.GetHashString()] = tx3
	txMap[tx5.GetHashString()] = tx5
	txMap[tx2.GetHashString()] = tx2
	txMap[tx6.GetHashString()] = tx6
	txMap[tx4.GetHashString()] = tx4
	txMap[tx1.GetHashString()] = tx1

	id := "2025 대선"

	pb := NewPreparedBlock(types.Topic(id), txMap)

	t.Logf(util.CyanString("PreparedBlock.VotingID: %s"), pb.VotingID)
	t.Logf(util.CyanString("PreparedBlock.MerkleRoot: %+v"), pb.MerkleRoot)
	t.Logf(util.CyanString("PreparedBlock.txx: %+v"), pb.txx)

	for i, tx := range pb.txx {
		t.Logf(util.CyanString("TX %+v : %s"), i+1, tx.GetHashString())
	}

	mr := "bb4a7eb0c49e083afe8e98aa2965291204aa5bfcc9d48ab21e8d8ed4168321bc"
	mrh, err := types.HashFromHashString(mr)

	assert.Nil(t, err)
	assert.Equal(t, mrh, pb.MerkleRoot)
	assert.Equal(t, mr, pb.MerkleRoot.String())

	t.Logf(util.MagentaString("MerkleRoot: %x"), pb.MerkleRoot)
	t.Logf(util.MagentaString("MerkleRootOracle: %x"), mrh)
	t.Logf(util.MagentaString("MerkleRootString: %s"), pb.MerkleRoot.String())
	t.Logf(util.MagentaString("MerkleRootStringOracle: %s"), mr)

	t.Log(util.GreenString("+-+-+-+-+-+-+-+-+- block_test.go::TestPreparedBlock \"END\" -+-+-+-+-+-+-+-+-+"))
}
