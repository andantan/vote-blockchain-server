package explorer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andantan/vote-blockchain-server/config"
	"github.com/andantan/vote-blockchain-server/core/blockchain"
	"github.com/andantan/vote-blockchain-server/core/mempool"
	"github.com/andantan/vote-blockchain-server/util"
)

type BlockChainExplorer struct {
	chain            *blockchain.BlockChain
	mempool          *mempool.MemPool
	baseDir          string
	blocksDir        string
	ExplorerPort     uint16
	ExplorerEndPoint string
}

func NewBlockChainExplorer(chain *blockchain.BlockChain, pool *mempool.MemPool) *BlockChainExplorer {
	systemBlockchainStoreBaseDir := config.GetEnvVar("SYSTEM_BLOCKCHAIN_STORE_BASE_DIR")
	systemBlockchainStoreBlockDir := config.GetEnvVar("SYSTEM_BLOCKCHAIN_STORE_BLOCK_DIR")
	connectionRestExplorerListenerPort := config.GetIntEnvVar("CONNECTION_REST_EXPLORER_LISTENER_PORT")
	connectionRestExplorerListenerEndpoint := config.GetEnvVar("CONNECTION_REST_EXPLORER_LISTENER_ENDPOINT")

	return &BlockChainExplorer{
		chain:            chain,
		mempool:          pool,
		baseDir:          systemBlockchainStoreBaseDir,
		blocksDir:        systemBlockchainStoreBlockDir,
		ExplorerPort:     uint16(connectionRestExplorerListenerPort),
		ExplorerEndPoint: connectionRestExplorerListenerEndpoint,
	}
}

func (e *BlockChainExplorer) Start() {
	http.HandleFunc(e.ExplorerEndPoint+"/block", e.handleBlockQuery)
	http.HandleFunc(e.ExplorerEndPoint+"/height", e.handleHeightQuery)
	http.HandleFunc(e.ExplorerEndPoint+"/headers", e.handleHeadersQuery)
	http.HandleFunc(e.ExplorerEndPoint+"/query", e.handleSpecQuery)
	http.HandleFunc(e.ExplorerEndPoint+"/mempool/pending", e.handleMempoolPendingsQuery)
	http.HandleFunc(e.ExplorerEndPoint+"/mempool/txx", e.handleMempoolTxxQuery)

	addr := fmt.Sprintf(":%d", e.ExplorerPort)

	log.Printf(util.SystemString("SYSTEM: Explorer listener opened { port: %d }"), e.ExplorerPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
