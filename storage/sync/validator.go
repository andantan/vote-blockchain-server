package sync

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/andantan/vote-blockchain-server/config"
	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/core/transaction"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
)

type SyncDir struct {
	baseDir   string
	blocksDir string
}

func NewSyncDir(baseDir, blocksDir string) SyncDir {
	return SyncDir{
		baseDir:   baseDir,
		blocksDir: blocksDir,
	}
}

type Validator struct {
	syncDir SyncDir
	headers []*block.Header
}

func NewValidator() *Validator {
	systemBlockchainStoreBaseDir := config.GetEnvVar("SYSTEM_BLOCKCHAIN_STORE_BASE_DIR")
	systemBlockchainStoreBlockDir := config.GetEnvVar("SYSTEM_BLOCKCHAIN_STORE_BLOCK_DIR")

	syncDir := NewSyncDir(systemBlockchainStoreBaseDir, systemBlockchainStoreBlockDir)

	return &Validator{
		syncDir: syncDir,
		headers: []*block.Header{},
	}
}

func (v *Validator) StartValidate() {
	log.Printf(
		util.CyanString("VALIDATE: Starting blockchain data validation. Data path: %s"),
		filepath.Join(v.syncDir.baseDir, v.syncDir.blocksDir),
	)

	startSyncTime := time.Now()

	bfs, err := v.getSortedBlockFilePaths()

	if err != nil {
		panic(err.Error())
	}

	if len(bfs) != 0 {
		blks, err := v.getSortedBlocks(bfs)

		if err != nil {
			panic(err.Error())
		}

		ph := types.ZeroHashCompact()

		for _, blk := range blks {
			if err := v.validateBlockByBlock(ph, blk); err != nil {
				panic(err.Error())
			}

			ph = blk.BlockHash
			v.syncBlockHeaders(blk)

			if blk.Height == 0 {
				log.Printf(util.CyanString("VALIDATE: Block( 0x%s ) with Height( %d ) => "+util.YellowString("Verified ( GENESIS BLOCK )")), blk.BlockHash.String(), blk.Height)
			}

			log.Printf(util.CyanString("VALIDATE: Block( 0x%s ) with Height( %d ) => "+util.YellowString("Verified")), blk.BlockHash.String(), blk.Height)
		}
	} else {
		log.Println(util.FatalString("SYSTEM: No existing blockchain data found. Initializing new blockchain with a genesis block."))
	}

	elapsedSyncTime := time.Since(startSyncTime)
	log.Printf(util.CyanString("VALIDATE: Blockchain validate data completed in %s"), elapsedSyncTime)
}

func (v *Validator) getSortedBlockFilePaths() ([]string, error) {
	return getJsonFiles(v.syncDir)
}

func (v *Validator) getSortedBlocks(blockFilePaths []string) ([]*block.Block, error) {
	blks := []*block.Block{}

	for _, file := range blockFilePaths {
		blk, err := loadBlock(file)

		if err != nil {
			return []*block.Block{}, err
		}

		blks = append(blks, blk)
	}

	return blks, nil
}

func (v *Validator) validateBlockByBlock(prevBlockHash types.Hash, blk *block.Block) error {
	if prevBlockHash != blk.PrevBlockHash {
		return fmt.Errorf(
			"chained block hash mismatch: block height %d (hash %s) previous hash does not match block height %d (hash %s)",
			blk.Height, blk.Header.PrevBlockHash.String(),
			blk.Height-1, prevBlockHash.String(),
		)
	}

	stx := transaction.NewSortedTxxFromJson(blk.Transactions)
	calculatedMerkleRoot := block.CalculateMerkleRoot(stx)

	if blk.MerkleRoot != calculatedMerkleRoot {
		return fmt.Errorf("block %d merkle root mismatch: stored %s, calculated %s",
			blk.Height,
			blk.Header.MerkleRoot.String(),
			calculatedMerkleRoot.String(),
		)
	}

	calculatedBlockHash := blk.Header.Hash()

	if blk.BlockHash != calculatedBlockHash {
		return fmt.Errorf("block %d hash mismatch: stored %s, calculated %s",
			blk.Height,
			blk.BlockHash.String(),
			calculatedBlockHash.String(),
		)
	}

	return nil
}

func (v *Validator) syncBlockHeaders(blk *block.Block) {
	v.headers = append(v.headers, blk.Header)
}

func (v *Validator) GetSyncedBlockHeaders() []*block.Header {
	return v.headers
}
