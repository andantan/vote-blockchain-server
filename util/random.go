package util

import (
	"fmt"
	"math/rand"

	"github.com/andantan/vote-blockchain-server/types"
)

func RandRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomBytes(size int) []byte {
	ticket := make([]byte, size)

	rand.Read(ticket)

	return ticket
}

func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(types.DIGEST_SIZE))
}

type OriginalVote struct {
	VoteHash   types.Hash
	UserHash   types.Hash
	VoteOption string
	Age        uint8
	Gender     rune
	VoteId     types.Proposal
}

func RandomVote() *OriginalVote {
	v := &OriginalVote{
		UserHash:   RandomHash(),
		VoteOption: "5",
		Age:        26,
		Gender:     'M',
		VoteId:     "2025-보건복지여론조사",
	}

	// fmt.Println(v.UserHash.String())

	data := fmt.Sprintf("%s|%s|%d|%c|%s",
		v.UserHash.String(), v.VoteOption, v.Age, v.Gender, v.VoteId)

	// fmt.Println(data)

	v.VoteHash = types.HashFromString(data)

	// fmt.Println(hex.EncodeToString(v.VoteHash.ToSlice()))

	return v
}
