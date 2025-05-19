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

type Vote struct {
	VoteHash   types.Hash
	UserHash   types.Hash
	VoteOption string
	Age        uint8
	Gender     rune
	ElectionId string
}

func RandomVote() *Vote {
	v := &Vote{
		VoteHash:   RandomHash(),
		UserHash:   RandomHash(),
		VoteOption: "5",
		Age:        26,
		Gender:     'M',
		ElectionId: "2025-보건복지여론조사",
	}

	return v
}
