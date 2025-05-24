package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashStringToBytes(t *testing.T) {
	str := "example string"
	hashedBytes := HashFromString(str)

	b, err := hex.DecodeString(hashedBytes.String())

	assert.Nil(t, err)

	for i := range 32 {
		assert.Equal(t, hashedBytes[i], b[i])
	}
}
