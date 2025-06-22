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
	_cfg_path := config.GetExplorerFilePathConfiguration()
	_cfg_end_point := config.GetExplorerListenerConfiguration()

	return &BlockChainExplorer{
		baseDir:          _cfg_path.ExplorerBaseDir,
		blocksDir:        _cfg_path.ExplorerBlockDir,
		ExplorerPort:     _cfg_end_point.ExplorerListenerPort,
		ExplorerEndPoint: _cfg_end_point.ExplorerListenerEndPoint,
	}
}

func (e *BlockChainExplorer) Start() {
	http.HandleFunc(e.ExplorerEndPoint, e.handleBlockQuery)

	addr := fmt.Sprintf(":%d", e.ExplorerPort)

	log.Printf(util.SystemString("SYSTEM: Explorer listener opened { port: %d }"), e.ExplorerPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
