package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func parse[T any](fileName string, cfg *T) {
	path := filepath.Join("./", "config", fileName)
	configFile, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("%s - reading error: %v", fileName, err)
	}
	if err = json.Unmarshal(configFile, cfg); err != nil {
		log.Fatalf("JSON unmarshalling failed: %v", err)
	}
}

const (
	SYSTEM_VALIDATOR_CONFIG_JSON              = "system_validator_config.json"
	SYSTEM_STORER_CONFIG_JSON                 = "system_storer_config.json"
	SYSTEM_EXPLORER_PATH_CONFIG_JSON          = "system_explorer_path_config.json"
	SYSTEM_CHAIN_PARAMS_CONFIG_JSON           = "system_chain_params_config.json"
	SYSTEM_CHANNEL_BUFFER_SIZE_CONFIG_JSON    = "system_channel_buffer_size_config.json"
	SYSTEM_PENDING_INTERNAL_TIMER_CONFIG_JSON = "system_pending_internal_timer_config.json"
)

const (
	CONNECTION_GRPC_VOTE_PROPOSAL_LISTENER_CONFIG_JSON = "connection_grpc_vote_proposal_listener_config.json"
	CONNECTION_GRPC_VOTE_SUBMIT_LISTENER_CONFIG_JSON   = "connection_grpc_vote_submit_listener_config.json"
	CONNECTION_REST_EXPLORER_LISTENER_CONFIG_JSON      = "connection_rest_explorer_listener_config.json"
	CONNECTION_UNICAST_PENDING_EVENT_CONFIG_JSON       = "connection_unicast_pending_event_config.json"
	CONNECTION_UNICAST_BLOCK_EVENT_CONFIG_JSON         = "connection_unicast_block_event_config.json"
)

func GetChannelBufferSizeSystemConfiguration() ChannelBufferSizeSystemConfiguration {
	cfgFileName := SYSTEM_CHANNEL_BUFFER_SIZE_CONFIG_JSON
	cfg := ChannelBufferSizeSystemConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetValidatorConfiguration() ValidatorConfiguration {
	cfgFileName := SYSTEM_VALIDATOR_CONFIG_JSON
	cfg := ValidatorConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetStorerConfiguration() StorerConfiguration {
	cfgFileName := SYSTEM_STORER_CONFIG_JSON
	cfg := StorerConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetExplorerFilePathConfiguration() ExplorerFilePathConfiguration {
	cfgFileName := SYSTEM_EXPLORER_PATH_CONFIG_JSON
	cfg := ExplorerFilePathConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetGrpcVoteProposalListenerConfiguration() GrpcVoteProposalListenerConfiguration {
	cfgFileName := CONNECTION_GRPC_VOTE_PROPOSAL_LISTENER_CONFIG_JSON
	cfg := GrpcVoteProposalListenerConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetGrpcVoteSubmitListenerConfiguration() GrpcVoteSubmitListenerConfiguration {
	cfgFileName := CONNECTION_GRPC_VOTE_SUBMIT_LISTENER_CONFIG_JSON
	cfg := GrpcVoteSubmitListenerConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetExplorerListenerConfiguration() ExplorerListenerConfiguration {
	cfgFileName := CONNECTION_REST_EXPLORER_LISTENER_CONFIG_JSON
	cfg := ExplorerListenerConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetChainParameterConfiguration() ChainParameterConfiguration {
	cfgFileName := SYSTEM_CHAIN_PARAMS_CONFIG_JSON
	cfg := ChainParameterConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetPendingEventUnicastConfiguration() PendingEventUnicastConfiguration {
	cfgFileName := CONNECTION_UNICAST_PENDING_EVENT_CONFIG_JSON
	cfg := PendingEventUnicastConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetBlockEventUnicastConfiguration() BlockEventUnicastConfiguration {
	cfgFileName := CONNECTION_UNICAST_BLOCK_EVENT_CONFIG_JSON
	cfg := BlockEventUnicastConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}

func GetPendingInternalTimerConfiguration() PendingInternalTimerConfiguration {
	cfgFileName := SYSTEM_PENDING_INTERNAL_TIMER_CONFIG_JSON
	cfg := PendingInternalTimerConfiguration{}

	parse(cfgFileName, &cfg)

	return cfg
}
