package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/andantan/vote-blockchain-server/config"
	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/storage/path"
	"github.com/andantan/vote-blockchain-server/util"
)

type JsonStorer struct {
	baseDir     string
	blocksDir   string
	blockSaveCh chan *block.Block
	wg          sync.WaitGroup
}

func NewStore() *JsonStorer {
	systemBlockchainStoreBaseDir := config.GetEnvVar("SYSTEM_BLOCKCHAIN_STORE_BASE_DIR")
	systemBlockchainStoreBlockDir := config.GetEnvVar("SYSTEM_BLOCKCHAIN_STORE_BLOCK_DIR")
	systemBlockStoreChannelBufferSize := config.GetIntEnvVar("SYSTEM_BLOCK_STORE_CHANNEL_BUFFER_SIZE")

	storer := &JsonStorer{}
	storer.setBaseDirectory(systemBlockchainStoreBaseDir, systemBlockchainStoreBlockDir)
	storer.setChannel(uint16(systemBlockStoreChannelBufferSize))
	storer.wg.Add(1)

	go storer.saveBlocks()

	return storer
}

func (js *JsonStorer) setBaseDirectory(baseDir, blocksDir string) {
	log.Println(util.SystemString("STORER: JsonStorer setting base directory..."))

	js.baseDir = baseDir
	js.blocksDir = blocksDir

	log.Println(util.SystemString("SYSTEM: JsonStorer base directory setting is done"))
}

func (js *JsonStorer) setChannel(bufferSize uint16) {
	log.Println(util.SystemString("SYSTEM: JsonStorer setting channel..."))

	js.blockSaveCh = make(chan *block.Block, bufferSize)

	log.Println(util.SystemString("SYSTEM: JsonStorer blockSave channel setting is done"))
}

func (js *JsonStorer) saveBlocks() {
	defer js.wg.Done()

	log.Println(util.StorerString("STORER: Starting blockSave receiver and processor goroutine"))

	blocksPath := filepath.Join(js.baseDir, js.blocksDir)

	if err := path.EnsureDir(blocksPath); err != nil {
		log.Fatalf(
			util.FatalString("STORER: Failed to create or verify block storage directory (%s): %v"),
			blocksPath, err,
		)
	}

	log.Printf(util.StorerString("STORER: Block storage directory '%s' ready"), blocksPath)

	for block := range js.blockSaveCh {

		jsonData, err := json.MarshalIndent(block, "", "  ")

		if err != nil {
			log.Printf(
				util.StorerString("STORER: Block %s (Height %d) Failed to marshalling: %v"),
				block.Header.VotingID,
				block.Header.Height,
				err,
			)
			continue
		}

		fileName := fmt.Sprintf("block_%d.json", block.Header.Height)
		filePath := filepath.Join(blocksPath, fileName)

		if err = os.WriteFile(filePath, jsonData, 0644); err != nil {
			log.Printf(
				util.FatalString("STORE: %s | 'block_%d.json' write failed: %v"),
				block.Header.VotingID,
				block.Header.Height,
				err,
			)

			continue
		}

		log.Printf(
			util.StorerString("STORER: %s | Successfully saved to file 'block_%d.json'"),
			block.Header.VotingID,
			block.Header.Height,
		)
	}
}

func (js *JsonStorer) SaveBlock(block *block.Block) {
	select {
	case js.blockSaveCh <- block:
		log.Printf(
			util.StorerString("STORER: %s | Block height %d successfully pushed to save channel"),
			block.Header.VotingID,
			block.Header.Height,
		)
	default:
		log.Printf(
			util.StorerString("STORER: Block save channel is full. Dropping block %s (height %d)"),
			block.Header.VotingID,
			block.Header.Height,
		)
	}
}

func (js *JsonStorer) Shutdown() {
	log.Println(util.StorerString("STORER: Initiating shutdown for JsonStorer. Closing blockSave channel"))

	close(js.blockSaveCh)
	js.wg.Wait()

	log.Println(util.StorerString("STORER: JsonStorer shutdown complete"))
}
