package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/core/transaction"
)

const (
	MESSAGE = `
--------------------------------------------------------------------------------------
| *H.Voting ID     : %-80s
| *H.Merkle Root   : %-80s
| *H.Height        : %-80d
| *H.PrevBlockHash : %-80s
| B.BlockHash      : %-80s
| B.TxLength       : %-80d
--------------------------------------------------------------------------------------
`
)
const (
	STORE_BASE_DIR   = "../../"
	STORE_BLOCKS_DIR = "blocks"
)

func main() {
	loadChain()
}

func loadChain() {
	blocksPath := filepath.Join(STORE_BASE_DIR, STORE_BLOCKS_DIR)

	if err := EnsureDir(blocksPath); err != nil {
		log.Fatalf(
			"STORER: Failed to create or verify block storage directory (%s): %v",
			blocksPath, err,
		)
	}

	files, err := GetFilesInDir(blocksPath, "block_*.json")

	if err != nil {
		log.Fatalf("GetFilesInDir error: %s", err.Error())
	}

	s, _ := SortBlockFilesByHeight(files)

	for _, file := range s {
		load(file)
	}
}

func load(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("block file not found: %s", path)
		}
		fmt.Printf("error reading block file %s: %s", path, err.Error())
	}

	blk := &block.Block{}
	err = json.Unmarshal(data, blk)

	if err != nil {
		fmt.Printf("error unmarshalling block data from %s: %s", path, err.Error())
	}

	log.Printf(MESSAGE, blk.VotingID, blk.MerkleRoot.String(), blk.Height, blk.PrevBlockHash.String(), blk.BlockHash, len(blk.Transactions))

	stx := transaction.NewSortedTxxFromJson(blk.Transactions)

	merkleRoot := block.CalculateMerkleRoot(stx)

	vmr := merkleRoot == blk.MerkleRoot

	log.Printf("MERKLEROOT VALIDATION: %t\n", vmr)

	vbh := blk.Header.Hash() == blk.BlockHash
	log.Printf("BLOCKHASH VALIDATION: %t\n", vbh)
	log.Printf("BLOCKHASH JSON: %s\n", blk.Header.Hash())
	log.Printf("BLOCKHASH: %s\n", blk.BlockHash)
}

func EnsureDir(path string) error {
	dir := filepath.Dir(path)

	log.Printf("dir: %s", dir)

	if dir == "." {
		dir = path
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}

	return nil
}

func GetFilesInDir(dirPath, pattern string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var matchedFiles []string

	parts := strings.Split(pattern, "*")
	prefix := parts[0]
	suffix := ""

	if len(parts) > 1 {
		suffix = parts[1]
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if strings.HasPrefix(fileName, prefix) && strings.HasSuffix(fileName, suffix) {
				matchedFiles = append(matchedFiles, filepath.Join(dirPath, fileName))
			}
		}
	}

	return matchedFiles, nil
}

type BlockFile struct {
	Path   string
	Height uint64
}

func SortBlockFilesByHeight(filePaths []string) ([]string, error) {
	var blockFiles []BlockFile

	for _, filePath := range filePaths {
		fileName := filepath.Base(filePath)

		if !strings.HasPrefix(fileName, "block_") || !strings.HasSuffix(fileName, ".json") {
			fmt.Printf("Warning: Skipping file with invalid name format: %s\n", fileName)
			continue
		}

		sHeight := strings.TrimSuffix(strings.TrimPrefix(fileName, "block_"), ".json")
		height, err := strconv.ParseUint(sHeight, 10, 64)

		if err != nil {
			return nil, fmt.Errorf("failed to parse block height from file name '%s': %w", fileName, err)
		}

		blockFiles = append(blockFiles, BlockFile{
			Path:   filePath,
			Height: height,
		})
	}

	sort.Slice(blockFiles, func(i, j int) bool {
		return blockFiles[i].Height < blockFiles[j].Height
	})

	sortedPaths := make([]string, len(blockFiles))

	for i, bf := range blockFiles {
		sortedPaths[i] = bf.Path
	}

	return sortedPaths, nil
}
