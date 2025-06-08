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
	"github.com/andantan/vote-blockchain-server/impulse-client/util"
)

type ProposalResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func BurstProposalClient(max int) {
	client := NewProposalClient(max)

	for i, vote := range client.Topics.Votes {
		if max <= i {
			break
		}

		go client.RequestProposal(vote)
	}

	client.Wg.Wait()
}

type ProposalClient struct {
	Client                 *http.Client
	Wg                     *sync.WaitGroup
	Topics                 config.Topics
	EndPoint               config.VoteProposalEndPoint
	MinimumRangeBurstClock int
	MaximumRangeBurstClock int
}

func NewProposalClient(max int) *ProposalClient {
	cfg := config.GetRequestBurstRangeClock()

	c := &ProposalClient{
		Client:                 &http.Client{Timeout: 10 * time.Second},
		Wg:                     &sync.WaitGroup{},
		Topics:                 config.GetTopics(),
		EndPoint:               config.GetVoteProposalEndPoint(),
		MinimumRangeBurstClock: int(cfg.RestProposalRequestsRandomMinimumSeconds),
		MaximumRangeBurstClock: int(cfg.RestProposalRequestsRandomMaximumSeconds),
	}

	if len(c.Topics.Votes) < max {
		panic(fmt.Sprintf("Cannot process %d proposals: only %d proposals available. 'max' value must not exceed the total number of proposals.", max, len(c.Topics.Votes)))
	}

	c.Wg.Add(max)

	return c
}

func (c *ProposalClient) GetUrl() string {
	return fmt.Sprintf("%s://%s:%d%s",
		c.EndPoint.RestVoteProposalProtocol,
		c.EndPoint.RestVoteProposalAddress,
		c.EndPoint.RestVoteProposalPort,
		c.EndPoint.RestVoteProposalEndPoint,
	)
}

func (c *ProposalClient) RequestProposal(vote config.Vote) {
	defer c.Wg.Done()

	time.Sleep(time.Duration(util.RandRange(c.MinimumRangeBurstClock, c.MaximumRangeBurstClock)) * time.Second)

	jsonData, err := json.Marshal(vote)

	if err != nil {
		log.Fatalf("error marshalling JSON: %v", err)
	}

	resp, err := c.Client.Post(
		c.GetUrl(),
		c.EndPoint.RestVoteProposalContentType,
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

	response := &ProposalResponse{}

	if err := json.Unmarshal(body, response); err != nil {
		log.Fatalf("error unmarshalling response JSON: %v", err)
	}

	log.Printf("Vote proposal response: %+v", response)
}
