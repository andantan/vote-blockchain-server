package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const (
	VOTE_DATA_JSON = "vote_data.json"
)

const (
	CONNECTION_REST_VOTE_PROPOSAL_CONFIG_JSON = "connection_rest_vote_proposal_config.json"
	CONNECTION_REST_VOTE_SUBMIT_CONFIG_JSON   = "connection_rest_vote_submit_config.json"
)

func parse[T any](fileName string, cfg *T) {
	path := filepath.Join("./", "data", fileName)
	configFile, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("%s - reading error: %v", fileName, err)
	}
	if err = json.Unmarshal(configFile, cfg); err != nil {
		log.Fatalf("JSON unmarshalling failed: %v", err)
	}
}

func GetTopics() Topics {
	cfgFileName := VOTE_DATA_JSON
	cfg := Topics{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetVoteProposalEndPoint() VoteProposalEndPoint {
	cfgFileName := CONNECTION_REST_VOTE_PROPOSAL_CONFIG_JSON
	cfg := VoteProposalEndPoint{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetVoteSubmitEndPoint() VoteSubmitEndPoint {
	cfgFileName := CONNECTION_REST_VOTE_SUBMIT_CONFIG_JSON
	cfg := VoteSubmitEndPoint{}

	parse(cfgFileName, &cfg)

	return cfg
}
