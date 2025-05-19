package core

import "github.com/andantan/vote-blockchain-server/types"

type Transaction struct {
	VoteHash   types.Hash
	VoteOption string
}
