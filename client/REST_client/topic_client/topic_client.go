package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	PROTOCOL = "http"
	ADDRESS  = "localhost"
	PORT     = 8080
	API      = "topic/new"
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
	url := fmt.Sprintf("%s://%s:%d/%s", PROTOCOL, ADDRESS, PORT, API)

	topic := NewTopic("2025 경선", 1)

	jsonData, err := json.Marshal(topic)

	if err != nil {
		log.Fatalf("error marshalling JSON: %v", err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))

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

	fmt.Printf("POST Response Status: %s\n", resp.Status)
	fmt.Printf("POST Response Body (unmarshalled): %+v", response)
}
