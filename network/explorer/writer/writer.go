package writer

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andantan/vote-blockchain-server/core/block"
	werror "github.com/andantan/vote-blockchain-server/error"
	"github.com/andantan/vote-blockchain-server/util"
)

func WriteJSONErrorResponse(w http.ResponseWriter, statusCode int, werror *werror.WrappedError) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(statusCode)

	log.Printf(util.CyanString("EXPLORER: Error [%s]: %s"), werror.Code, werror.Message)

	jsonResponse := ExplorerBlockAPIResponse{
		Success: "false",
		Message: werror.Error(),
		Status:  werror.Code,
	}

	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Printf(util.RedString("EXPLORER: Failed to write JSON error response: %v"), err)
	}
}

func WriteJSONSuccessBlockResponse(w http.ResponseWriter, blk *block.Block) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)

	jsonResponse := ExplorerBlockAPIResponse{
		Success: "true",
		Message: "Operation successful",
		Status:  "OK",
		Block:   blk,
	}

	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Printf(util.RedString("EXPLORER: Failed to write JSON success response: %v"), err)
	}
}
