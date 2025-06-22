package data

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const (
	REGISTER_DATA_JSON = "register_data.json"
	USER_DATA_JSON     = "user_data.json"
	VOTE_DATA_JSON     = "vote_data.json"
	BALLOT_DATA_JSON   = "ballot_options_data.json"
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

func GetRegisterData() Topics {
	cfgFileName := REGISTER_DATA_JSON
	cfg := Topics{}

	parse(cfgFileName, &cfg)

	return cfg
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

func GetUsers() Users {
	cfgFileName := USER_DATA_JSON
	cfg := Users{}

	parse(cfgFileName, &cfg)

	return cfg
}
