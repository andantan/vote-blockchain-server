package transaction

import "github.com/andantan/vote-blockchain-server/types"

type Transaction struct {
	Hash   types.Hash
	Option string
}
