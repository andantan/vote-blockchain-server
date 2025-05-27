package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	DIGEST_SIZE = 32 // SHA-256 output size 32Byte
)

// [ Tx, Header, ... ] -> SHA-256 -> 32Byte
type Hash [DIGEST_SIZE]uint8

func EmptyHash() Hash {
	return ZeroHashCompact()
}

func NilHash() Hash {
	return FFHashCompact()
}

func FilledHash(b byte) Hash {
	h := Hash{}

	for i := range DIGEST_SIZE {
		h[i] = b
	}

	return h
}

// Equal with FilledHash(0x00) or Hash{}
func ZeroHashCompact() Hash {
	return Hash{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}

// Equal with FilledHash(0xFF)
func FFHashCompact() Hash {
	return Hash{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
}

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

func HashFromString(s string) Hash {
	return sha256.Sum256([]byte(s))
}

// b []byte len must be 32 & Return Hash that casted from bytes
func HashFromBytes(b []byte) Hash {
	if len(b) != DIGEST_SIZE {
		msg := fmt.Sprintf("given bytes with length %d should be 32", len(b))

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

// If string is valid then return Hash, true
// else return nil, false
func IsValidHashString(s string) (Hash, bool) {
	h, err := hex.DecodeString(s)

	if err != nil {
		return Hash{}, false
	}

	return Hash(h), true
}

func HashFromHashString(s string) (Hash, error) {
	h, err := hex.DecodeString(s)

	if err != nil {
		return Hash{}, err
	}

	return Hash(h), nil
}
