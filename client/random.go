package main

import "crypto/rand"

Hash = [32]uint8

type Vote struct {
	VoteHash   [32]uint8
	UserHash   [32]uint8
	VoteOption string
	Age        uint8
	Gender     rune
	ElectionId string
}

func RandomVote() *Vote {
	v := &Vote{
		VoteHash:   types.RandomHash(),
		UserHash:   types.RandomHash(),
		VoteOption: "5",
		Age:        26,
		Gender:     'M',
		ElectionId: "2025-보건복지여론조사사",
	}

	return v
}

func RandomBytes(size int) []byte {
	ticket := make([]byte, size)

	rand.Read(ticket)

	return ticket
}

func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(types.DIGEST_SIZE))
}
