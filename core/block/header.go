package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"github.com/andantan/vote-blockchain-server/types"
)

// Height 0 will be a genesis block
type Header struct {
	VotingID      types.Topic `json:"voting_id"`
	MerkleRoot    types.Hash  `json:"merkle_root"`
	Height        uint64      `json:"height"`
	PrevBlockHash types.Hash  `json:"prev_block_hash"`
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.Write([]byte(h.VotingID))
	buf.Write(h.MerkleRoot.ToSlice())

	binary.Write(buf, binary.BigEndian, h.Height)

	buf.Write(h.PrevBlockHash.ToSlice())

	return buf.Bytes()
}

func (h *Header) Hash() types.Hash {
	hash := sha256.Sum256(h.Bytes())

	return types.Hash(hash)
}
