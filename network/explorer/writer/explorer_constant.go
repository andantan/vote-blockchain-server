package writer

import "github.com/andantan/vote-blockchain-server/core/block"

type ExplorerBlockAPIResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`

	Block *block.Block `json:"block"`
}

type ExplorerHeightAPIResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`

	Height uint32 `json:"height"`
}

type ResponseHeader struct {
	VotingID      string `json:"voting_id"`
	Height        uint64 `json:"height"`
	MerkleRoot    string `json:"merkle_root"`
	BlockHash     string `json:"block_hash"`
	PrevBlockHash string `json:"prev_block_hash"`
}

func NewResponseHeader(h *block.Header) *ResponseHeader {
	return &ResponseHeader{
		VotingID:      string(h.VotingID),
		Height:        h.Height,
		MerkleRoot:    "0x" + h.MerkleRoot.String(),
		BlockHash:     "0x" + h.Hash().String(),
		PrevBlockHash: "0x" + h.PrevBlockHash.String(),
	}
}

type ExplorerHeadersAPIResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`

	From    uint32            `json:"from"`
	To      uint32            `json:"to"`
	Headers []*ResponseHeader `json:"headers"`
}

type ExplorerSpecAPIResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`

	Type string            `json:"type"`
	Spec []*ResponseHeader `json:"headers"`
}
