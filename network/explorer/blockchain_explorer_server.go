package explorer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andantan/vote-blockchain-server/config"
	"github.com/andantan/vote-blockchain-server/util"
)

type BlockChainExplorer struct {
	baseDir          string
	blocksDir        string
	ExplorerPort     uint16
	ExplorerEndPoint string
}

func NewBlockChainExplorer() *BlockChainExplorer {
	systemBlockchainStoreBaseDir := config.GetEnvVar("SYSTEM_BLOCKCHAIN_STORE_BASE_DIR")
	systemBlockchainStoreBlockDir := config.GetEnvVar("SYSTEM_BLOCKCHAIN_STORE_BLOCK_DIR")
	connectionRestExplorerListenerPort := config.GetIntEnvVar("CONNECTION_REST_EXPLORER_LISTENER_PORT")
	connectionRestExplorerListenerEndpoint := config.GetEnvVar("CONNECTION_REST_EXPLORER_LISTENER_ENDPOINT")

	return &BlockChainExplorer{
		baseDir:          systemBlockchainStoreBaseDir,
		blocksDir:        systemBlockchainStoreBlockDir,
		ExplorerPort:     uint16(connectionRestExplorerListenerPort),
		ExplorerEndPoint: connectionRestExplorerListenerEndpoint,
	}
}

func (e *BlockChainExplorer) Start() {
	http.HandleFunc(e.ExplorerEndPoint, e.handleBlockQuery)

	addr := fmt.Sprintf(":%d", e.ExplorerPort)

	log.Printf(util.SystemString("SYSTEM: Explorer listener opened { port: %d }"), e.ExplorerPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
