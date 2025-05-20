package mempool

import "github.com/andantan/vote-blockchain-server/types"

type MemPool interface {
}

type PendingPool struct {
	pendings map[types.ElectionID]*Pending
}
