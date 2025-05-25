package block

import (
	"crypto/sha256"

	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
)

func CalculateMerkleRoot(txx *transaction.SortedTxx) types.Hash {
	if txx == nil {
		return types.NilHash()
	}

	if txx.Len() == 0 {
		return types.EmptyHash()
	}

	hashes := txx.GetHashSlice()

	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	for len(hashes) > 1 {
		nextLevelHashes := make([]types.Hash, 0, len(hashes)/2)

		for i := 0; i < len(hashes); i += 2 {
			combinedHashBytes := append(hashes[i].ToSlice(), hashes[i+1].ToSlice()...)

			newHash := sha256.Sum256(combinedHashBytes)
			nextLevelHashes = append(nextLevelHashes, types.Hash(newHash))
		}
		hashes = nextLevelHashes

		if len(hashes)%2 != 0 && len(hashes) > 1 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}
	}

	return hashes[0]
}
