package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/impulse-client/config"
	"github.com/andantan/vote-blockchain-server/impulse-client/data"
	"github.com/andantan/vote-blockchain-server/impulse-client/util"
)

type SubmitRequest struct {
	Hash   string `json:"hash"`
	Option string `json:"option"`
	Topic  string `json:"topic"`
}

func NewSubmitRequest(hash, option, topic string) *SubmitRequest {
	return &SubmitRequest{
		Hash:   hash,
		Option: option,
		Topic:  topic,
	}
}

type SubmitResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	Status         string `json:"status"`
	Topic          string `json:"topic"`
	HttpStatusCude int    `json:"http_status_code"`
	UserHash       string `json:"user_hash"`
	VoteHash       string `json:"vote_hash"`
	VoteOption     string `json:"vote_option"`
}

func BurstSubmitClient(max int, votes []data.Vote) {
	client := NewSubmitClient(max, votes)

	for i, vote := range client.Topics {
		if max <= i {
			break
		}

		go client.RequestSubmitLoop(vote)
	}

	client.Wg.Wait()
}

type SubmitClient struct {
	Client                 *http.Client
	Wg                     *sync.WaitGroup
	Topics                 []data.Vote
	EndPoint               config.VoteSubmitEndPoint
	MinimumRangeBurstClock int
	MaximumRangeBurstClock int
}

func NewSubmitClient(max int, votes []data.Vote) *SubmitClient {
	cfg := config.GetRequestBurstRangeClock()

	c := &SubmitClient{
		Client:                 &http.Client{Timeout: 10 * time.Second},
		Wg:                     &sync.WaitGroup{},
		Topics:                 votes,
		EndPoint:               config.GetVoteSubmitEndPoint(),
		MinimumRangeBurstClock: int(cfg.RestSubmitRequestsRandomMinimunMilliSeconds),
		MaximumRangeBurstClock: int(cfg.RestSubmitRequestsRandomMaximumMilliSeconds),
	}

	if len(c.Topics) < max {
		panic(fmt.Sprintf("Cannot process %d proposals: only %d proposals available. 'max' value must not exceed the total number of proposals.", max, len(c.Topics)))
	}

	c.Wg.Add(max)

	return c
}

func (c *SubmitClient) GetUrl() string {
	return fmt.Sprintf("%s://%s:%d%s",
		c.EndPoint.RestVoteSubmitProtocol,
		c.EndPoint.RestVoteSubmitAddress,
		c.EndPoint.RestVoteSubmitPort,
		c.EndPoint.RestVoteSubmitEndPoint,
	)
}

func (c *SubmitClient) RequestSubmitLoop(vote data.Vote) {
	defer c.Wg.Done()

	log.Printf("RequestLoop %.20s start", vote.Topic)

	requestCount := 0
	requestOption := make(map[string]int)

	ballotOptions := data.GetBallotOptions()
	users := data.GetUsers()

	users.ShuffleUserHashs()

	for i := range users.UserLength {
		randOpt := util.RandOption(ballotOptions.BallotOptions)

		v := NewSubmitRequest(
			users.GetUserHash(i),
			randOpt,
			vote.Topic,
		)

		response := c.RequestSubmit(v)

		if !response.Success {
			break
		}

		requestCount++
		requestOption[randOpt]++

		time.Sleep(time.Duration(util.RandRange(c.MinimumRangeBurstClock, c.MaximumRangeBurstClock)) * time.Millisecond)
	}

	log.Printf("RequestLoop %s exit | { requestCount: %d , result: %+v }", vote.Topic, requestCount, requestOption)
}

func (c *SubmitClient) RequestSubmit(vote *SubmitRequest) *SubmitResponse {
	jsonData, err := json.Marshal(vote)

	if err != nil {
		log.Fatalf("error marshalling JSON: %v", err)
	}

	resp, err := c.Client.Post(
		c.GetUrl(),
		c.EndPoint.RestVoteSubmitContentType,
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		log.Fatalf("error POST request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	response := &SubmitResponse{}

	if err := json.Unmarshal(body, response); err != nil {
		log.Fatalf("error unmarshalling response JSON: %v", err)
	}

	log.Printf("Vote submit Topic: %s, response: %+v", vote.Topic, response)

	return response
}
