package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/andantan/vote-blockchain-server/util"
)

const (
	PROTOCOL = "http"
	ADDRESS  = "localhost"
	PORT     = 8080
	API      = "vote/submit"
)

var URL string = fmt.Sprintf("%s://%s:%d/%s", PROTOCOL, ADDRESS, PORT, API)

type Vote struct {
	Hash   string `json:"hash"`
	Option string `json:"option"`
	Topic  string `json:"topic"`
}

func NewVote(hash, option, topic string) *Vote {
	return &Vote{
		Hash:   hash,
		Option: option,
		Topic:  topic,
	}
}

type VoteResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {

	topics := []string{
		"2025 대선",
		"2025 경선",
		"2025 보건의료 여론조사",
		"법률개정안 찬반 투표",
		"상법개정안 시범 기간 조사",
	}

	blockTimer := time.NewTicker(2 * time.Minute)
	defer blockTimer.Stop()

	for _, topic := range topics {
		go RequestLoop(topic)
	}

	<-blockTimer.C

}

func RequestVote(vote *Vote) *VoteResponse {
	jsonData, err := json.Marshal(vote)

	if err != nil {
		log.Fatalf("error marshalling JSON: %v", err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Post(URL, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalf("error POST request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	response := VoteResponse{}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("error unmarshalling response JSON: %v", err)
	}

	return &response
}

func RequestLoop(topic string) {
	for {
		vote := NewVote(
			util.RandomHash().String(),
			randOpt(),
			topic,
		)

		response := RequestVote(vote)

		if strings.Compare(response.Success, "false") == 0 {
			break
		}

		log.Printf(util.YellowString("Response: { %+v }"), response)

		time.Sleep(time.Duration(util.RandRange(40, 200)) * time.Millisecond)
	}

	log.Printf(util.CyanString("RequestLoop %s exit"), topic)
}

func randOpt() string {
	options := []rune("12345")

	return string(options[rand.Intn(len(options))])
}
