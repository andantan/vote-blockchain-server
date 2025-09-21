package explorer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/andantan/vote-blockchain-server/core/block"
	werror "github.com/andantan/vote-blockchain-server/error"
	"github.com/andantan/vote-blockchain-server/network/explorer/writer"
	"github.com/andantan/vote-blockchain-server/types"
	"github.com/andantan/vote-blockchain-server/util"
)

func (e *BlockChainExplorer) handleBlockQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorMessage := fmt.Sprintf("%s method not allowed", r.Method)
		wrappedError := werror.NewWrappedError("METHOD_NOT_ALLOWED", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusMethodNotAllowed, wrappedError)

		return
	}

	queriedHeightStr := r.URL.Query().Get("height")

	if queriedHeightStr == "" {
		errorMessage := "Query parameter 'height' is required."
		wrappedError := werror.NewWrappedError("EMPTY_QUERY_PARAMETER", errorMessage, nil)

		// log.Printf(util.CyanString("EXPLORER: Bad request: %s"), errorMessage)
		writer.WriteJSONErrorResponse(w, http.StatusBadRequest, wrappedError)

		return
	}

	height, err := strconv.Atoi(queriedHeightStr)

	if err != nil || height < 0 {
		errorMessage := fmt.Sprintf("Invalid 'height' value: %s. Must be a non-negative integer.", queriedHeightStr)
		wrappedError := werror.NewWrappedError("INVALID_QUERY_PARAMETER", errorMessage, nil)

		// log.Printf(util.CyanString("EXPLORER: Bad request: %s"), errorMessage)
		writer.WriteJSONErrorResponse(w, http.StatusBadRequest, wrappedError)

		return
	}

	fileName := fmt.Sprintf("block_%d.json", height)
	filePath := filepath.Join(e.baseDir, e.blocksDir, fileName)
	blockData, err := os.ReadFile(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			errorMessage := fmt.Sprintf("Block at height %d not found.", height)
			wrappedError := werror.NewWrappedError("BLOCK_NOT_FOUND", errorMessage, nil)

			writer.WriteJSONErrorResponse(w, http.StatusNotFound, wrappedError)

			return
		}

		errorMessage := fmt.Sprintf("Failed to read block data for height %d.", height)
		wrappedError := werror.NewWrappedError("BLOCK_READ_ERROR", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusNotFound, wrappedError)

		return
	}

	blk := new(block.Block)

	if err := json.Unmarshal(blockData, blk); err != nil {
		errorMessage := fmt.Sprintf("error unmarshalling block data for height %d", height)
		wrappedError := werror.NewWrappedError("UNMARSHALLING_ERROR", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusInternalServerError, wrappedError)

		return
	}

	writer.WriteJSONSuccessBlockResponse(w, blk)

	log.Printf(util.CyanString("EXPLORER: Successfully served block query %d from %s"), height, filePath)
}

func (e *BlockChainExplorer) handleHeightQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorMessage := fmt.Sprintf("%s method not allowed", r.Method)
		wrappedError := werror.NewWrappedError("METHOD_NOT_ALLOWED", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusMethodNotAllowed, wrappedError)

		return
	}

	height := e.chain.Height()
	writer.WriteJSONSuccessHeightResponse(w, height)

	log.Printf(util.CyanString("EXPLORER: Successfully served height(%d) query"), height)
}

func (e *BlockChainExplorer) handleHeadersQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorMessage := fmt.Sprintf("%s method not allowed", r.Method)
		wrappedError := werror.NewWrappedError("METHOD_NOT_ALLOWED", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusMethodNotAllowed, wrappedError)

		return
	}

	queriedFromStr := r.URL.Query().Get("from")

	if queriedFromStr == "" {
		errorMessage := "Query parameter 'from' is required."
		wrappedError := werror.NewWrappedError("EMPTY_QUERY_PARAMETER", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusBadRequest, wrappedError)

		return
	}

	from, err := strconv.Atoi(queriedFromStr)

	if err != nil || from < 0 {
		errorMessage := fmt.Sprintf("Invalid 'from' value: %s. Must be a non-negative integer.", queriedFromStr)
		wrappedError := werror.NewWrappedError("INVALID_QUERY_PARAMETER", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusBadRequest, wrappedError)

		return
	}

	queriedToStr := r.URL.Query().Get("to")

	if queriedToStr == "" {
		errorMessage := "Query parameter 'to' is required."
		wrappedError := werror.NewWrappedError("EMPTY_QUERY_PARAMETER", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusBadRequest, wrappedError)

		return
	}

	to, err := strconv.Atoi(queriedToStr)

	if err != nil || to < 0 {
		errorMessage := fmt.Sprintf("Invalid 'to' value: %s. Must be a non-negative integer.", queriedToStr)
		wrappedError := werror.NewWrappedError("INVALID_QUERY_PARAMETER", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusBadRequest, wrappedError)

		return
	}

	headers := e.chain.GetHeadersByRange(uint32(from), uint32(to))

	writer.WriteJSONSuccessHeadersResponse(w, uint32(from), uint32(to), headers)

	log.Printf(util.CyanString("EXPLORER: Successfully served headers from=%d, to=%d query"), from, to)
}

func (e *BlockChainExplorer) handleSpecQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorMessage := fmt.Sprintf("%s method not allowed", r.Method)
		wrappedError := werror.NewWrappedError("METHOD_NOT_ALLOWED", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusMethodNotAllowed, wrappedError)

		return
	}

	targetFromStr := r.URL.Query().Get("target")

	if targetFromStr == "" {
		errorMessage := "Query parameter 'target' is required."
		wrappedError := werror.NewWrappedError("EMPTY_QUERY_PARAMETER", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusBadRequest, wrappedError)

		return
	}

	targetFromStr = strings.TrimPrefix(targetFromStr, "0x")
	headers := e.chain.GetHeadersByRange(0, e.chain.Height())
	resHeaders := make([]*block.Header, 0)

	var queryTarget string
	for _, h := range headers {
		if targetFromStr == string(h.VotingID) {
			queryTarget = "id"
			resHeaders = append(resHeaders, h)
		}

		if targetFromStr == h.Proposer.String() {
			queryTarget = "proposer"
			resHeaders = append(resHeaders, h)
		}

		if targetFromStr == h.Hash().String() {
			queryTarget = "block_hash"
			resHeaders = append(resHeaders, h)
			break
		}

		if targetFromStr == h.MerkleRoot.String() {
			queryTarget = "merkle_root"
			resHeaders = append(resHeaders, h)
			break
		}
	}

	if queryTarget == "" {
		queryTarget = "null"
	}

	writer.WriteJSONSuccessSpecResponse(w, resHeaders, queryTarget)

	log.Printf(util.CyanString("EXPLORER: Successfully served spec types=%s, len=%d query"), queryTarget, len(resHeaders))
}

func (e *BlockChainExplorer) handleMempoolPendingsQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorMessage := fmt.Sprintf("%s method not allowed", r.Method)
		wrappedError := werror.NewWrappedError("METHOD_NOT_ALLOWED", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusMethodNotAllowed, wrappedError)

		return
	}

	pendings := e.mempool.SeekPendings()
	res := make([]writer.ResponsePending, 0)

	for proposal, pending := range *pendings {
		p := writer.ResponsePending{
			Proposal: string(proposal),
			Proposer: "0x" + pending.GetProposer(),
			Option:   pending.OptCache,
		}

		res = append(res, p)
	}

	writer.WriteJSONSuccessPendingsResponse(w, res)

	log.Printf(util.CyanString("EXPLORER: Successfully served pending len=%d query"), len(res))
}

func (e *BlockChainExplorer) handleMempoolTxxQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorMessage := fmt.Sprintf("%s method not allowed", r.Method)
		wrappedError := werror.NewWrappedError("METHOD_NOT_ALLOWED", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusMethodNotAllowed, wrappedError)

		return
	}

	pendings := e.mempool.SeekPendings()
	idFromStr := r.URL.Query().Get("id")

	if idFromStr == "" {
		errorMessage := "Query parameter 'id' is required."
		wrappedError := werror.NewWrappedError("EMPTY_QUERY_PARAMETER", errorMessage, nil)

		writer.WriteJSONErrorResponse(w, http.StatusBadRequest, wrappedError)

		return
	}

	pool := make(map[string]string)

	if !e.mempool.IsOpen(types.Proposal(idFromStr)) {
		res := writer.ResponseTxx{
			Proposal: "null",
			Proposer: "0x" + types.ZeroHashCompact().String(),
			Pool:     pool,
		}

		writer.WriteJSONSuccessTxxResponse(w, res)

		return
	}

	p := (*pendings)[types.Proposal(idFromStr)]
	res := writer.ResponseTxx{
		Proposal: idFromStr,
		Proposer: "0x" + p.GetProposer(),
		Pool:     make(map[string]string),
	}

	for _, tx := range p.Txx {
		res.Pool["0x"+tx.GetHashString()] = tx.Option
	}

	writer.WriteJSONSuccessTxxResponse(w, res)

	log.Printf(util.CyanString("EXPLORER: Successfully served pending id=%s len=%d query"), idFromStr, len(res.Pool))

}
