package data

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const (
	VOTE_DATA_JSON   = "vote_data.json"
	BALLOT_DATA_JSON = "ballot_options_data.json"
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

func GetBallotOptions() BallotOption {
	cfgFileName := BALLOT_DATA_JSON
	cfg := BallotOption{}

	parse(cfgFileName, &cfg)

	return cfg
}
