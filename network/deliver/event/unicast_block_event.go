package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/util"
)

type CreatedBlockEventEndpoint struct {
	*DefaultEndPoint
}

func NewCreatedBlockEventEndpoint(protocol, address string, port uint16, path string) *CreatedBlockEventEndpoint {
	return &CreatedBlockEventEndpoint{
		&DefaultEndPoint{
			protocol: protocol,
			address:  address,
			port:     port,
			path:     path,
		},
	}
}

func (e *CreatedBlockEventEndpoint) getUrl() string {
	return fmt.Sprintf("%s://%s:%d%s", e.protocol, e.address, e.port, e.path)
}

type CreatedBlockEventRequest struct {
	VoteId string `json:"vote_id"`
	Height uint64 `json:"height"`
}

func NewCreatedBlockEventRequest(voteId string, height uint64) *CreatedBlockEventRequest {
	return &CreatedBlockEventRequest{
		VoteId: voteId,
		Height: height,
	}
}

type CreatedBlockEventResponse struct {
	Caching bool   `json:"caching"`
	Message string `json:"message"`
	VoteId  string `json:"vote_id"`
	Height  uint64 `json:"height"`
}

type CreatedBlockeventUnicaster struct {
	endPoint *CreatedBlockEventEndpoint
	client   *http.Client
}

func NewCreatedBlockeventUnicaster(cfg *CreatedBlockEventEndpoint) *CreatedBlockeventUnicaster {
	return &CreatedBlockeventUnicaster{
		endPoint: cfg,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (u *CreatedBlockeventUnicaster) GetUrl() string {
	return u.endPoint.getUrl()
}

func (u *CreatedBlockeventUnicaster) Unicast(createdBlock *block.Block) {
	req := NewCreatedBlockEventRequest(string(createdBlock.VotingID), createdBlock.Height)

	buf, _ := json.Marshal(req)

	log.Printf(
		util.DeliverString("DELIVER: BlockCreatedEventUnicaster.Unicast request { vote_id: %s, height: %d }"),
		req.VoteId, req.Height,
	)

	res, err := u.client.Post(u.endPoint.getUrl(), JSON_CONTENT_TYPE, bytes.NewBuffer(buf))

	if err != nil {
		log.Printf(util.FatalString("error BlockCreatedEventUnicaster POST request: %v"), err)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	dataReq := &CreatedBlockEventResponse{}

	if err := json.Unmarshal(body, dataReq); err != nil {
		log.Printf(util.FatalString("error BlockCreatedEventUnicaster unmarshalling response JSON: %v"), err)
	}

	log.Printf(
		util.DeliverString("DELIVER: BlockCreatedEventUnicaster.Unicast response { vote_id: %s, height: %d, caching: %t, message: %s }"),
		dataReq.VoteId, dataReq.Height, dataReq.Caching, dataReq.Message,
	)
}
