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
	t.Log(util.TestDecoratorString("+-+-+-+-+-+-+-+-+- block_test.go::TestHeaderBytes \"RUN\" -+-+-+-+-+-+-+-+-+"))

	randomMerkleRootHash := util.RandomHash()
	randomPrevBlockHash := util.RandomHash()

	t.Logf(util.TestInfoString("randomMerkleRootHash %v"), randomMerkleRootHash)
	t.Logf(util.TestInfoString("randomPrevBlockHash %v"), randomPrevBlockHash)

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

	t.Log(util.TestDecoratorString("+-+-+-+-+-+-+-+-+- block_test.go::TestHeaderBytes \"END\" -+-+-+-+-+-+-+-+-+"))
}

func TestProtoBlock(t *testing.T) {
	t.Log(util.TestDecoratorString("+-+-+-+-+-+-+-+-+- block_test.go::TestProtoBlock \"RUN\" -+-+-+-+-+-+-+-+-+"))

	pb := getProtoBlockForTest(t)

	t.Logf(util.TestInfoString("PreparedBlock.VotingID: %s"), pb.VotingID)
	t.Logf(util.TestInfoString("PreparedBlock.MerkleRoot: %+v"), pb.MerkleRoot)
	t.Logf(util.TestInfoString("PreparedBlock.txx: %+v"), pb.txx)

	for i, tx := range pb.txx {
		t.Logf(util.TestInfoString("TX %+v : %s"), i+1, tx.GetHashString())
	}

	mr := "a34caf1a5723242da60fa57ad5f85623ec6ec84721ecbd289e23655c95779dea"
	mrh, err := types.HashFromHashString(mr)

	assert.Nil(t, err)
	assert.Equal(t, mrh, pb.MerkleRoot)
	assert.Equal(t, mr, pb.MerkleRoot.String())

	t.Logf(util.TestOracleString("MerkleRoot: %x"), pb.MerkleRoot)
	t.Logf(util.TestOracleString("MerkleRootOracle: %x"), mrh)
	t.Logf(util.TestOracleString("MerkleRootString: %s"), pb.MerkleRoot.String())
	t.Logf(util.TestOracleString("MerkleRootStringOracle: %s"), mr)

	t.Log(util.TestDecoratorString("+-+-+-+-+-+-+-+-+- block_test.go::TestProtoBLock \"END\" -+-+-+-+-+-+-+-+-+"))
}

func TestCreateBlock(t *testing.T) {
	t.Log(util.TestDecoratorString("+-+-+-+-+-+-+-+-+- block_test.go::TestCreateBlock \"RUN\" -+-+-+-+-+-+-+-+-+"))

	gb := GenesisBlock()

	t.Logf(util.TestInfoString("GenesisBlock -> *Header.VotingID: %s"), gb.Header.VotingID)
	t.Logf(util.TestInfoString("GenesisBlock -> *Header.MerkleRoot: %s"), gb.Header.MerkleRoot)
	t.Logf(util.TestInfoString("GenesisBlock -> *Header.Height: %d"), gb.Header.Height)
	t.Logf(util.TestInfoString("GenesisBlock -> *Header.PrevBlockHash: %s"), gb.Header.PrevBlockHash)
	t.Logf(util.TestInfoString("GenesisBlock -> BlockHash: %s"), gb.BlockHash)
	t.Logf(util.TestInfoString("GenesisBlock -> len(Txx): %d"), len(gb.Transactions))

	genesisBlockHashString := "67df818365d4af91b8f47434118287781eebf39e49290aff9f3d9909ddc7a9c2"
	genesisBlockHashOracle, err := types.HashFromHashString(genesisBlockHashString)

	assert.Nil(t, err)
	assert.Equal(t, genesisBlockHashOracle, gb.BlockHash)

	t.Logf(util.TestOracleString("GenesisBlockHash: %s"), gb.BlockHash)
	t.Logf(util.TestOracleString("GenesisBlockHashOracle: %s"), genesisBlockHashOracle)

	pb := getProtoBlockForTest(t)

	block := NewBlockFromPrevHeader(gb.Header, pb)

	t.Logf(util.TestInfoString("CreatedBlock -> *Header.VotingID: %s"), block.Header.VotingID)
	t.Logf(util.TestInfoString("CreatedBlock -> *Header.MerkleRoot: %s"), block.Header.MerkleRoot)
	t.Logf(util.TestInfoString("CreatedBlock -> *Header.Height: %d"), block.Header.Height)
	t.Logf(util.TestInfoString("CreatedBlock -> *Header.PrevBlockHash: %s"), block.Header.PrevBlockHash)
	t.Logf(util.TestInfoString("CreatedBlock -> BlockHash: %s"), block.BlockHash)
	t.Logf(util.TestInfoString("CreatedBlock -> len(Txx): %d"), len(block.Transactions))

	createdBlockHashString := "1427c493814cbb169cc116c012eefc8ab7fc454390345b9d9dc9212c5edcae9a"
	createdBlockHashOracle, err := types.HashFromHashString(createdBlockHashString)

	assert.Nil(t, err)
	assert.Equal(t, block.PrevBlockHash, gb.BlockHash)
	assert.Equal(t, createdBlockHashOracle, block.BlockHash)

	t.Logf(util.TestOracleString("GenesisBlock.BlockHash: %s"), gb.BlockHash)
	t.Logf(util.TestOracleString("CreatedBlock.PrevBlockHash: %s"), block.PrevBlockHash)
	t.Logf(util.TestOracleString("CreatedBlockHash: %s"), block.BlockHash)
	t.Logf(util.TestOracleString("CreatedBlockHashOracle: %s"), createdBlockHashOracle)

	t.Log(util.TestDecoratorString("+-+-+-+-+-+-+-+-+- block_test.go::TestCreateBlock \"END\" -+-+-+-+-+-+-+-+-+"))
}

func getProtoBlockForTest(t *testing.T) *ProtoBlock {
	h1, err1 := types.HashFromHashString("3d2f41f8505ba7310ce68debc73618e7f20a3a656027125918b30701f13f6b4c")
	h2, err2 := types.HashFromHashString("c4d03172e0021fe3a66cad7a07e78e75d66338ee43c12733c352d1e12015ee2f")
	h3, err3 := types.HashFromHashString("d48e038e70c74fbc5d7dbef37acae71c7fffce6ae566cc518563bc2a459e6ec3")
	h4, err4 := types.HashFromHashString("732dc5537d95ff76294c3d1a3666dd3eb8038ce0eaf4982806b416d4bbdd085f")
	h5, err5 := types.HashFromHashString("9b28eb0ca548229990b85ac566799a4a569a2d928f10807ffe349b06ae13beec")
	h6, err6 := types.HashFromHashString("c8b74d3493891e65858ce9a0b903275bee03d5ecd6437d17cac4d43d057ef6d3")
	h7, err7 := types.HashFromHashString("6474a5df58aa06ba3906248b8d18e91b34ba3bd801f1d18a2a5b01033b1be0f0")
	h8, err8 := types.HashFromHashString("8d24d306fdb3fd057c8123286dcbf89e74e3ba8cd9ddc8cda2e8cae2c0810daf")
	h9, err9 := types.HashFromHashString("d937b2c6a94b7378621d15166fb633bbead0e641af04535041336df1c9fd0e36")
	h10, err10 := types.HashFromHashString("d8eec4e6d2b316e60f27f8a118b2a3759552fbc2ab215d8e178e829b8dba3c5a")
	h11, err11 := types.HashFromHashString("e485446df84ee0ccb35c2e7cb0a4e1a6f8e5d8f39695929bd770a1e7ac0c5e3c")
	h12, err12 := types.HashFromHashString("f06ed42542ebb62ff3baf1726f3d26959eb7ca678b4c06c46b41e648dbd561dc")

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.Nil(t, err4)
	assert.Nil(t, err5)
	assert.Nil(t, err6)
	assert.Nil(t, err7)
	assert.Nil(t, err8)
	assert.Nil(t, err9)
	assert.Nil(t, err10)
	assert.Nil(t, err11)
	assert.Nil(t, err12)

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
	time.Sleep(10 * time.Millisecond)
	tx7 := transaction.NewTransaction(h7, "2", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx8 := transaction.NewTransaction(h8, "1", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx9 := transaction.NewTransaction(h9, "2", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx10 := transaction.NewTransaction(h10, "3", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx11 := transaction.NewTransaction(h11, "3", time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	tx12 := transaction.NewTransaction(h12, "1", time.Now().UnixNano())

	txMap := make(map[string]*transaction.Transaction)

	txMap[tx3.GetHashString()] = tx3
	txMap[tx12.GetHashString()] = tx12
	txMap[tx9.GetHashString()] = tx9
	txMap[tx10.GetHashString()] = tx10
	txMap[tx5.GetHashString()] = tx5
	txMap[tx1.GetHashString()] = tx1
	txMap[tx11.GetHashString()] = tx11
	txMap[tx2.GetHashString()] = tx2
	txMap[tx6.GetHashString()] = tx6
	txMap[tx8.GetHashString()] = tx8
	txMap[tx4.GetHashString()] = tx4
	txMap[tx7.GetHashString()] = tx7

	id := "2025 대선"

	return NewProtoBlock(types.Proposal(id), txMap)
}
