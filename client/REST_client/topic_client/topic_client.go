package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/andantan/vote-blockchain-server/util"
)

const (
	PROTOCOL = "http"
	ADDRESS  = "localhost"
	PORT     = 8080
	API      = "topic/new"
)

var (
	URL string = fmt.Sprintf("%s://%s:%d/%s", PROTOCOL, ADDRESS, PORT, API)
)

type Topic struct {
	Topic    string `json:"topic"`
	Duration int    `json:"duration"`
}

func NewTopic(topic string, duration int) *Topic {
	return &Topic{
		Topic:    topic,
		Duration: duration,
	}
}

type TopicResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {

	topics := []*Topic{
		NewTopic("2025 대선", 1),
		NewTopic("2025 경선", 2),
		NewTopic("2025 보건의료 여론조사", 2),
		NewTopic("법률개정안 찬반 투표", 1),
		NewTopic("상법개정안 시범 기간 조사", 2),
	}

	for _, topic := range topics {
		response := RequestTopic(topic)

		fmt.Printf("POST Response Body: %+v\n", response)

		time.Sleep(time.Duration(util.RandRange(1, 3)) * time.Second)
	}
}

func RequestTopic(topic *Topic) *TopicResponse {
	jsonData, err := json.Marshal(topic)

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

	response := TopicResponse{}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("error unmarshalling response JSON: %v", err)
	}

	return &response
}
