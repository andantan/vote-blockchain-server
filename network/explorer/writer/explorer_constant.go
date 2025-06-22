package writer

import "github.com/andantan/vote-blockchain-server/core/block"

type ExplorerBlockAPIResponse struct {
	Success string       `json:"success"`
	Message string       `json:"message"`
	Status  string       `json:"status"`
	Block   *block.Block `json:"block"`
}
