package util

import (
	"crypto/rand"

	"github.com/andantan/vote-blockchain-server/types"
)

func RandomBytes(size int) []byte {
	ticket := make([]byte, size)

	rand.Read(ticket)

	return ticket
}

func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(types.DIGEST_SIZE))
}
