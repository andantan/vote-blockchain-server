package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"github.com/andantan/vote-blockchain-server/types"
)

// Height 0 will be a genesis block
type Header struct {
	VotingID      types.Topic
	MerkleRoot    types.Hash // Hashs txx
	Height        uint64
	PrevBlockHash types.Hash // Chaining with HeaderHash
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
