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

type VoteProposalRequest struct {
	Topic    string `json:"topic"`
	Duration int    `json:"duration"`
}

func NewVoteProposalRequest(topic string, duration int) *VoteProposalRequest {
	return &VoteProposalRequest{
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

	proposals := []*VoteProposalRequest{
		NewVoteProposalRequest("2025 대선", 1),
		NewVoteProposalRequest("2025 경선", 2),
		NewVoteProposalRequest("2025 보건의료 여론조사", 3),
		NewVoteProposalRequest("법률개정안 찬반 투표", 4),
		NewVoteProposalRequest("상법개정안 시범 기간 조사", 5),
		NewVoteProposalRequest("기후변화 대응 방안 선호도 조사", 6),
		NewVoteProposalRequest("인공지능 교육 도입 찬반 설문", 7),
		NewVoteProposalRequest("수원시 대중교통 만족도 평가", 8),
		NewVoteProposalRequest("청년 주거 정책 의견 수렴", 9),
		NewVoteProposalRequest("국민연금 개편안 대국민 토론", 10),
		NewVoteProposalRequest("미래 식량 기술 투자 필요성 조사", 1),
		NewVoteProposalRequest("문화예술 바우처 사업 확대 여부", 2),
		NewVoteProposalRequest("자율주행 자동차 상용화 시점 예측", 3),
		NewVoteProposalRequest("코로나19 재유행 대비 행동 지침", 4),
		NewVoteProposalRequest("초고령사회 대비 사회복지 시스템 개선", 5),
		NewVoteProposalRequest("2026 지방선거 후보자 적합도", 6),
		NewVoteProposalRequest("공공 의료시설 확충 필요성", 7),
		NewVoteProposalRequest("초등학생 코딩 교육 의무화", 8),
		NewVoteProposalRequest("친환경 에너지 전환 정책 평가", 9),
		NewVoteProposalRequest("디지털 교과서 도입 만족도", 10),
		NewVoteProposalRequest("MZ세대 공정성 인식 조사", 1),
		NewVoteProposalRequest("수도권 주택 공급 확대 방안", 2),
		NewVoteProposalRequest("소상공인 지원 정책 효과 분석", 3),
		NewVoteProposalRequest("청소년 마약 예방 교육 강화", 4),
		NewVoteProposalRequest("노인 일자리 창출 정책 개선", 5),
		NewVoteProposalRequest("동물 복지법 강화 찬반", 6),
		NewVoteProposalRequest("온라인 플랫폼 규제 방안", 7),
		NewVoteProposalRequest("농어촌 지역 소멸 위기 대응", 8),
		NewVoteProposalRequest("국가 채무 증가에 대한 우려", 9),
		NewVoteProposalRequest("K-콘텐츠 해외 진출 전략", 10),
	}

	max := 14

	wg.Add(max)

	for i, topic := range proposals {
		if max <= i {
			break
		}

		go RequestTopic(topic, wg)
	}

	wg.Wait()
}

func RequestTopic(v *VoteProposalRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Duration(util.RandRange(2, 6)) * time.Second)

	jsonData, err := json.Marshal(v)

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
