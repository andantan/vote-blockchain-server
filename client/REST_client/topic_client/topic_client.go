package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
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
	wg := &sync.WaitGroup{}

	topics := []*Topic{
		NewTopic("2025 대선", 2),
		NewTopic("2025 경선", 1),
		NewTopic("2025 보건의료 여론조사", 2),
		NewTopic("법률개정안 찬반 투표", 3),
		NewTopic("상법개정안 시범 기간 조사", 3),
		NewTopic("기후변화 대응 방안 선호도 조사", 1),
		NewTopic("인공지능 교육 도입 찬반 설문", 1),
		NewTopic("수원시 대중교통 만족도 평가", 3),
		NewTopic("청년 주거 정책 의견 수렴", 1),
		NewTopic("국민연금 개편안 대국민 토론", 3),
		NewTopic("미래 식량 기술 투자 필요성 조사", 2),
		NewTopic("문화예술 바우처 사업 확대 여부", 2),
		NewTopic("자율주행 자동차 상용화 시점 예측", 1),
		NewTopic("코로나19 재유행 대비 행동 지침", 1),
		NewTopic("초고령사회 대비 사회복지 시스템 개선", 3),
		NewTopic("2026 지방선거 후보자 적합도", 1),
		NewTopic("공공 의료시설 확충 필요성", 2),
		NewTopic("초등학생 코딩 교육 의무화", 3),
		NewTopic("친환경 에너지 전환 정책 평가", 2),
		NewTopic("디지털 교과서 도입 만족도", 3),
		NewTopic("MZ세대 공정성 인식 조사", 1),
		NewTopic("수도권 주택 공급 확대 방안", 1),
		NewTopic("소상공인 지원 정책 효과 분석", 3),
		NewTopic("청소년 마약 예방 교육 강화", 2),
		NewTopic("노인 일자리 창출 정책 개선", 4),
		NewTopic("동물 복지법 강화 찬반", 3),
		NewTopic("온라인 플랫폼 규제 방안", 3),
		NewTopic("농어촌 지역 소멸 위기 대응", 4),
		NewTopic("국가 채무 증가에 대한 우려", 1),
		NewTopic("K-콘텐츠 해외 진출 전략", 2),
	}

	max := 6

	wg.Add(max)

	for i, topic := range topics {
		if max <= i {
			break
		}

		go RequestTopic(topic, wg)
		time.Sleep(200 * time.Millisecond)
	}

	wg.Wait()
}

func RequestTopic(topic *Topic, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Duration(util.RandRange(2, 6)) * time.Second)

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

	log.Printf(util.CyanString("POST Response Body: %+v"), response)
}
