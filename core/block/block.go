package block

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

type Header struct {
	VotingID      types.Topic
	MerkleRoot    types.Hash // Hashs of all of transaction
	PrevBlockHash types.Hash // Chaining with HeaderHash
	Height        uint64
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

func CreateNewBlock() {
	fmt.Println("Created new block")
}
