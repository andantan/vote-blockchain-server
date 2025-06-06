package config

// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// System layer - Channel buffer size configuration

type ChannelBufferSizeSystemConfiguration struct {
	GrpcVoteProposalChannelBufferSize   uint16 `json:"GrpcVoteProposalChannelBufferSize"`
	GrpcVoteSubmitChannelBufferSize     uint16 `json:"GrpcVoteSubmitChannelBufferSize"`
	PendingTransactionChannelBufferSize uint16 `json:"PendingTransactionChannelBufferSize"`
	PendedPropaginateChannelBufferSize  uint16 `json:"PendedPropaginateChannelBufferSize"`
	BlockPropaginatedChannelBufferSize  uint16 `json:"BlockPropaginatedChannelBufferSize"`
}

// #end
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// Block storage path

type ValidatorConfiguration struct {
	StoreBaseDir  string `json:"StoreBaseDir"`
	StoreBlockDir string `json:"StoreBlockDir"`
}

type StorerConfiguration struct {
	StoreBaseDir  string `json:"StoreBaseDir"`
	StoreBlockDir string `json:"StoreBlockDir"`
}

// #end
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// gRPC listener configuration

type GrpcVoteProposalListenerConfiguration struct {
	Network string `json:"ProposalGrpcListenerNetwork"`
	Port    uint16 `json:"ProposalGrpcListenerPort"`
}

type GrpcVoteSubmitListenerConfiguration struct {
	Network string `json:"SubmitGrpcListenerNetwork"`
	Port    uint16 `json:"SubmitGrpcListenerPort"`
}

// #end
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// Chain parameter ( mempool configuration )

type ChainParameterConfiguration struct {
	BlockIntervalSeconds uint32 `json:"BlockIntervalSeconds"`
	MaxTransactionSize   uint32 `json:"MaxTransactionSize"`
}

// #end
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// Event deliver ( Unicast ) configuration

type PendingEventUnicastConfiguration struct {
	PendingEventUnicastProtocol        string `json:"PendingEventUnicastProtocol"`
	PendingEventUnicastAddress         string `json:"PendingEventUnicastAddress"`
	PendingEventUnicastPort            uint16 `json:"PendingEventUnicastPort"`
	ExpiredPendingEventUnicastEndPoint string `json:"ExpiredPendingEventUnicastEndPoint"`
}

type BlockEventUnicastConfiguration struct {
	BlockEventUnicastProtocol        string `json:"BlockEventUnicastProtocol"`
	BlockEventUnicastAddress         string `json:"BlockEventUnicastAddress"`
	BlockEventUnicastPort            uint16 `json:"BlockEventUnicastPort"`
	CreatedBlockEventUnicastEndPoint string `json:"CreatedBlockEventUnicastEndPoint"`
}

// #end
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// Pending timer configuration

type PendingInternalTimerConfiguration struct {
	ResetTimeDurationSeconds      uint8 `json:"ResetTimeDurationSeconds"`
	InturruptTimerDurationSeconds uint8 `json:"InturruptTimerDurationSeconds"`
	CloseTimerDurationSeconds     uint8 `json:"CloseTimerDurationSeconds"`
}

// #end
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
