package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/andantan/vote-blockchain-server/util"
)

var httpClient *http.Client

const (
	PROTOCOL = "http"
	ADDRESS  = "localhost"
	PORT     = 8080
	API      = "/vote/submit"
)

var URL string = fmt.Sprintf("%s://%s:%d%s", PROTOCOL, ADDRESS, PORT, API)

type VoteSubmitRequest struct {
	Hash   string `json:"hash"`
	Option string `json:"option"`
	Topic  string `json:"topic"`
}

func NewVoteSubmitRequest(hash, option, topic string) *VoteSubmitRequest {
	return &VoteSubmitRequest{
		Hash:   hash,
		Option: option,
		Topic:  topic,
	}
}

type VoteSubmitResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	var max int
	flag.IntVar(&max, "max", 10, "Limit submit loop entity")
	flag.Parse()

	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	wg := &sync.WaitGroup{}

	proposals := []string{
		"2025 대선",
		"2025 경선",
		"2025 보건의료 여론조사",
		"법률개정안 찬반 투표",
		"상법개정안 시범 기간 조사",
		"기후변화 대응 방안 선호도 조사",
		"인공지능 교육 도입 찬반 설문",
		"수원시 대중교통 만족도 평가",
		"청년 주거 정책 의견 수렴",
		"국민연금 개편안 대국민 토론",
		"미래 식량 기술 투자 필요성 조사",
		"문화예술 바우처 사업 확대 여부",
		"자율주행 자동차 상용화 시점 예측",
		"코로나19 재유행 대비 행동 지침",
		"초고령사회 대비 사회복지 시스템 개선",
		"2026 지방선거 후보자 적합도",
		"공공 의료시설 확충 필요성",
		"초등학생 코딩 교육 의무화",
		"친환경 에너지 전환 정책 평가",
		"디지털 교과서 도입 만족도",
		"MZ세대 공정성 인식 조사",
		"수도권 주택 공급 확대 방안",
		"소상공인 지원 정책 효과 분석",
		"청소년 마약 예방 교육 강화",
		"노인 일자리 창출 정책 개선",
		"동물 복지법 강화 찬반",
		"온라인 플랫폼 규제 방안",
		"농어촌 지역 소멸 위기 대응",
		"국가 채무 증가에 대한 우려",
		"K-콘텐츠 해외 진출 전략",
	}

	if len(proposals) < max {
		panic(fmt.Sprintf("Cannot process %d submits: only %d proposals available. 'max' value must not exceed the total number of proposals.", max, len(proposals)))
	}

	wg.Add(max)

	for i, topic := range proposals {
		if max <= i {
			break
		}

		go RequestLoop(topic, wg)
	}

	wg.Wait()
}

func RequestVote(v *VoteSubmitRequest) *VoteSubmitResponse {
	jsonData, err := json.Marshal(v)

	if err != nil {
		log.Fatalf("error marshalling JSON: %v", err)
	}

	resp, err := httpClient.Post(URL, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalf("error POST request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	r := VoteSubmitResponse{}

	if err := json.Unmarshal(body, &r); err != nil {
		log.Fatalf("error unmarshalling response JSON: %v", err)
	}

	return &r
}

func RequestLoop(topic string, wg *sync.WaitGroup) {
	// time.Sleep(time.Duration(util.RandRange(5, 30)) * time.Second)

	defer wg.Done()

	log.Printf(util.YellowString("RequestLoop %.20s start"), topic)

	requestCount := 0
	requestOption := make(map[string]int)

	for {
		randOpt := randOpt()

		vote := NewVoteSubmitRequest(
			util.RandomHash().String(),
			randOpt,
			topic,
		)

		response := RequestVote(vote)

		if strings.Compare(response.Success, "false") == 0 {
			break
		}

		requestCount++
		requestOption[randOpt]++
		// log.Printf(util.YellowString("Response: { %+v }"), response)

		time.Sleep(time.Duration(util.RandRange(40, 60)) * time.Millisecond)
	}
	log.Printf(util.CyanString("RequestLoop %s exit | { requestCount: %d , result: %+v }"), topic, requestCount, requestOption)
}

func randOpt() string {
	options := []rune("12345")

	return string(options[rand.Intn(len(options))])
}
