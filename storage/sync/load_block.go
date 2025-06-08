package sync

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
	"github.com/andantan/vote-blockchain-server/storage/path"
)

func loadBlock(path string) (*block.Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("block file not found: %s", path)
		}

		return nil, fmt.Errorf("error reading block file %s: %s", path, err.Error())
	}

	blk := &block.Block{}
	err = json.Unmarshal(data, blk)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling block data from %s: %s", path, err.Error())
	}

	return blk, nil
}

type blockFile struct {
	path   string
	height uint64
}

func getJsonFiles(syncDir SyncDir) ([]string, error) {
	blocksPath := filepath.Join(syncDir.baseDir, syncDir.blocksDir)

	if err := path.EnsureDir(blocksPath); err != nil {
		log.Fatalf(
			"STORER: Failed to create or verify block storage directory (%s): %v",
			blocksPath, err,
		)
	}

	files, err := path.GetFilesInDir(blocksPath, "block_*.json")

	if err != nil {
		log.Fatalf("GetFilesInDir error: %s", err.Error())
	}

	return sortBlockFilesByHeight(files)
}

func sortBlockFilesByHeight(filePaths []string) ([]string, error) {
	var blockFiles []blockFile

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

		blockFiles = append(blockFiles, blockFile{
			path:   filePath,
			height: height,
		})
	}

	sort.Slice(blockFiles, func(i, j int) bool {
		return blockFiles[i].height < blockFiles[j].height
	})

	sortedPaths := make([]string, len(blockFiles))

	for i, bf := range blockFiles {
		sortedPaths[i] = bf.path
	}

	return sortedPaths, nil
}
