package util

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/rand"
)

func RandRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandOption(params []rune) string {
	return string(params[rand.Intn(len(params))])
}

func RandomHashString() string {
	return hex.EncodeToString(ToSlice(HashFromBytes(RandomBytes(32))))
}

func RandomBytes(size int) []byte {
	ticket := make([]byte, size)

	crand.Read(ticket)

	return ticket
}

func HashFromBytes(b []byte) [32]uint8 {
	var t [32]uint8

	for i := range 32 {
		t[i] = b[i]
	}

	return t
}

func ToSlice(d [32]uint8) []byte {
	b := make([]byte, 32)

	for i := range 32 {
		b[i] = d[i]
	}

	return b
}
