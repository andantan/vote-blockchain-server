package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type Header struct {
	VotingID      types.Topic
	MerkleRoot    types.Hash // Hashs of all of transaction
	Height        uint64
	PrevBlockHash types.Hash // Chaining with HeaderHash
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}

	gob.NewEncoder(buf).Encode(h)

	return buf.Bytes()
}

type Block struct {
	*Header
	HeaderHash   types.Hash // Hashs Header, header has merkleroot
	Transactions []*transaction.Transaction
}

type PreparedBlock struct {
	VotingID   types.Topic
	MerkleRoot types.Hash
	txx        []*transaction.Transaction
}

func NewPreparedBlock(ID types.Topic, txMap map[string]*transaction.Transaction) *PreparedBlock {
	stx := transaction.NewTxMapSorter(txMap)
	txHashes := GetHashSlice(stx)
	merkleRoot := CalculateMerkleRoot(txHashes)

	return &PreparedBlock{
		VotingID:   ID,
		MerkleRoot: merkleRoot,
		txx:        stx.GetTxx(),
	}
}

func GetHashSlice(txx *transaction.SortedTxx) []types.Hash {
	hashes := make([]types.Hash, 0, len(txx.GetTxx()))

	for _, tx := range txx.GetTxx() {
		hashes = append(hashes, tx.GetHash())
	}

	return hashes
}

func CalculateMerkleRoot(hashes []types.Hash) types.Hash {
	if len(hashes) == 0 {
		return types.Hash{}
	}

	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	for len(hashes) > 1 {
		nextLevelHashes := make([]types.Hash, 0, len(hashes)/2)

		for i := 0; i < len(hashes); i += 2 {
			combinedHashBytes := append(hashes[i].ToSlice(), hashes[i+1].ToSlice()...)

			newHash := sha256.Sum256(combinedHashBytes)
			nextLevelHashes = append(nextLevelHashes, types.Hash(newHash))
		}
		hashes = nextLevelHashes

		if len(hashes)%2 != 0 && len(hashes) > 1 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}
	}

	return hashes[0]
}

func GetHashStringSlice(txx *transaction.SortedTxx) []string {
	hashes := make([]string, 0, len(txx.GetTxx()))

	for _, tx := range txx.GetTxx() {
		hashes = append(hashes, tx.GetHashString())
	}

	return hashes
}

func CalculateMerkleRootFromStrings(hashes []string) string {
	if len(hashes) == 0 {
		return types.Hash{}.String()
	}

	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	for len(hashes) > 1 {
		nextLevelHashes := make([]string, 0, len(hashes)/2)

		for i := 0; i < len(hashes); i += 2 {
			combinedHashString := hashes[i] + hashes[i+1]

			newHashBytes := types.Hash(sha256.Sum256([]byte(combinedHashString)))

			nextLevelHashes = append(nextLevelHashes, newHashBytes.String())
		}

		hashes = nextLevelHashes

		if len(hashes)%2 != 0 && len(hashes) > 1 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}
	}

	return hashes[0]
}
