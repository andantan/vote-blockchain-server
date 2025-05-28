package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/storage/path"
	"github.com/andantan/vote-blockchain-server/util"
)

const (
	BLOCK_SAVE_CHANNEL_BUFFER_SIZE = 128
)

type JsonStorer struct {
	baseDir     string
	blocksDir   string
	blockSaveCh chan *block.Block
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewStore(baseDir, blocksDir string) *JsonStorer {
	ctx, cancel := context.WithCancel(context.Background())

	js := &JsonStorer{
		ctx:    ctx,
		cancel: cancel,
	}

	js.setBaseDirectory(baseDir, blocksDir)
	js.setChannel()
	js.wg.Add(1)

	go js.saveBlocks()

	return js
}

func (js *JsonStorer) setBaseDirectory(baseDir, blocksDir string) {
	log.Printf(util.SystemString("STORER: JsonStorer setting base directory { baseDir: %s, blocksDir: %s }"),
		baseDir, blocksDir)

	js.baseDir = baseDir
	js.blocksDir = blocksDir

	log.Println(util.SystemString("SYSTEM: JsonStorer base directory setting is done"))
}

func (js *JsonStorer) setChannel() {
	log.Printf(
		util.SystemString("SYSTEM: JsonStorer setting channel... | { BLOCK_SAVE_CHANNEL_BUFFER_SIZE: %d }"),
		BLOCK_SAVE_CHANNEL_BUFFER_SIZE,
	)

	js.blockSaveCh = make(chan *block.Block, BLOCK_SAVE_CHANNEL_BUFFER_SIZE)

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

	for {
		select {
		case block := <-js.blockSaveCh:

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

		case <-js.ctx.Done():
			return
		}
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
	js.cancel()
	js.wg.Wait()
	close(js.blockSaveCh)
	log.Println(util.StorerString("STORER: JsonStorer shutdown complete"))
}
