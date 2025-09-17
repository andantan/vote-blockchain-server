package writer

import "github.com/andantan/vote-blockchain-server/core/block"

type BasicResponseStatus struct {
}

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

type ExplorerHeadersAPIResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`

	From    uint32          `json:"from"`
	To      uint32          `json:"to"`
	Headers []*block.Header `json:"headers"`
}

type ExplorerSpecAPIResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`

	Type string          `json:"type"`
	Spec []*block.Header `json:"headers"`
}
