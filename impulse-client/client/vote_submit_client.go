package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func BurstSubmitClient(max int) {
	client := NewSubmitClient(max)

	for i, vote := range client.Topics.Votes {
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
	Topics                 data.Topics
	EndPoint               config.VoteSubmitEndPoint
	MinimumRangeBurstClock int
	MaximumRangeBurstClock int
}

func NewSubmitClient(max int) *SubmitClient {
	cfg := config.GetRequestBurstRangeClock()

	c := &SubmitClient{
		Client:                 &http.Client{Timeout: 10 * time.Second},
		Wg:                     &sync.WaitGroup{},
		Topics:                 data.GetTopics(),
		EndPoint:               config.GetVoteSubmitEndPoint(),
		MinimumRangeBurstClock: int(cfg.RestSubmitRequestsRandomMinimunMilliSeconds),
		MaximumRangeBurstClock: int(cfg.RestSubmitRequestsRandomMaximumMilliSeconds),
	}

	if len(c.Topics.Votes) < max {
		panic(fmt.Sprintf("Cannot process %d proposals: only %d proposals available. 'max' value must not exceed the total number of proposals.", max, len(c.Topics.Votes)))
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

	for {
		randOpt := util.RandOption(ballotOptions.BallotOptions)

		v := NewSubmitRequest(
			util.RandomHashString(),
			randOpt,
			vote.Topic,
		)

		response := c.RequestSubmit(v)

		if strings.Compare(response.Success, "false") == 0 {
			fmt.Printf("%+v\n", response)
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

	log.Printf("Vote submit response: %+v", response)

	return response
}
