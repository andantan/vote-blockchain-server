package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/andantan/vote-blockchain-server/core/mempool"
	"github.com/andantan/vote-blockchain-server/util"
)

type ExpiredPendingEventEndPoint struct {
	*DefaultEndPoint
}

func NewExpiredPendingEventEndPoint(protocol, address string, port uint16, path string) *ExpiredPendingEventEndPoint {
	return &ExpiredPendingEventEndPoint{
		&DefaultEndPoint{
			protocol: protocol,
			address:  address,
			port:     port,
			path:     path,
		},
	}
}

func (e *ExpiredPendingEventEndPoint) getUrl() string {
	return fmt.Sprintf("%s://%s:%d%s", e.protocol, e.address, e.port, e.path)
}

type ExpiredPendingEventRequest struct {
	VoteId      string         `json:"vote_id"`
	VoteCount   int            `json:"vote_count"`
	VoteOptions map[string]int `json:"vote_options"`
}

func NewExpiredPendingEventRequest(voteId string, voteCount int, voteOptions map[string]int) *ExpiredPendingEventRequest {
	return &ExpiredPendingEventRequest{
		VoteId:      voteId,
		VoteCount:   voteCount,
		VoteOptions: voteOptions,
	}
}

type ExpiredPendingEventResponse struct {
	Cached    bool   `json:"success"`
	Status    string `json:"status"`
	VoteId    string `json:"vote_id"`
	VoteCount int    `json:"vote_count"`
}

type ExpiredPendingEventUnicaster struct {
	endPoint *ExpiredPendingEventEndPoint
	client   *http.Client
}

func NewExpiredPendingEventUnicaster(cfg *ExpiredPendingEventEndPoint) *ExpiredPendingEventUnicaster {
	return &ExpiredPendingEventUnicaster{
		endPoint: cfg,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (u *ExpiredPendingEventUnicaster) GetUrl() string {
	return u.endPoint.getUrl()
}

func (u *ExpiredPendingEventUnicaster) Unicast(expiredPending *mempool.Pended) {
	req := NewExpiredPendingEventRequest(
		string(expiredPending.GetPendingID()),
		expiredPending.GetCachedLength(),
		expiredPending.GetCachedOptions(),
	)

	buf, _ := json.Marshal(req)

	res, err := u.client.Post(u.endPoint.getUrl(), JSON_CONTENT_TYPE, bytes.NewBuffer(buf))

	if err != nil {
		log.Printf(util.FatalString("error ExpiredPendingEventUnicaster POST request: %v"), err)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	dataReq := &ExpiredPendingEventResponse{}

	if err := json.Unmarshal(body, dataReq); err != nil {
		log.Printf(util.FatalString("error ExpiredPendingEventUnicaster unmarshalling response JSON: %v"), err)
	}

	log.Printf(
		util.DeliverString("DELIVER: ExpiredPendingEventUnicaster.Unicast response { voting_id: %s, count: %d, caching: %t, status: %s }"),
		dataReq.VoteId, dataReq.VoteCount, dataReq.Cached, dataReq.Status,
	)
}
