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
	VotingId      string         `json:"voting_id"`
	SubmitsLength int            `json:"submits_length"`
	SubmitsOption map[string]int `json:"submits_option"`
}

func NewExpiredPendingEventRequest(votingId string, submitsLength int, submitsOption map[string]int) *ExpiredPendingEventRequest {
	return &ExpiredPendingEventRequest{
		VotingId:      votingId,
		SubmitsLength: submitsLength,
		SubmitsOption: submitsOption,
	}
}

type ExpiredPendingEventResponse struct {
	Caching  bool   `json:"caching"`
	VotingID string `json:"voting_id"`
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
		util.DeliverString("DELIVER: ExpiredPendingEventUnicaster.Unicast response { voting_id: %s, caching: %t  }"),
		dataReq.VotingID, dataReq.Caching,
	)
}
