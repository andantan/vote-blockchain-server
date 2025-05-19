package types

import (
	"encoding/hex"
	"fmt"
)

const (
	DIGEST_SIZE = 32 // SHA-256 output size 32Byte
)

// [ Tx, Header, ... ] -> SHA-256 -> 32Byte
type Hash [DIGEST_SIZE]uint8

func (h Hash) IsZero() bool {
	for i := range DIGEST_SIZE {
		if h[i] != 0 {
			return false
		}
	}

	return true
}

func (h Hash) ToSlice() []byte {
	b := make([]byte, 32)

	for i := range DIGEST_SIZE {
		b[i] = h[i]
	}

	return b
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

// b []byte len must be 32 & Return Hash that casted from bytes
func HashFromBytes(b []byte) Hash {
	if len(b) != DIGEST_SIZE {
		msg := fmt.Sprintf("Given bytes with length %d should be 32", len(b))

		// System can not continue
		panic(msg)
	}

	var t [DIGEST_SIZE]uint8

	// Byte slice element by element. Not clone
	for i := range DIGEST_SIZE {
		t[i] = b[i]
	}

	return Hash(t)
}
