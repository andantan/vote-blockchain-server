package block

import (
	"testing"

	"github.com/andantan/vote-blockchain-server/util"
	"github.com/stretchr/testify/assert"
)

func TestHeaderBytes(t *testing.T) {
	randomMerkleRootHash := util.RandomHash()
	randomPrevBlockHash := util.RandomHash()

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
}
