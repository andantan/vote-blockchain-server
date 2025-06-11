package config

type VoteProposalEndPoint struct {
	RestVoteProposalProtocol    string `json:"RestVoteProposalProtocol"`
	RestVoteProposalAddress     string `json:"RestVoteProposalAddress"`
	RestVoteProposalPort        uint16 `json:"RestVoteProposalPort"`
	RestVoteProposalEndPoint    string `json:"RestVoteProposalEndPoint"`
	RestVoteProposalContentType string `json:"RestVoteProposalContentType"`
}

type VoteSubmitEndPoint struct {
	RestVoteSubmitProtocol    string `json:"RestVoteSubmitProtocol"`
	RestVoteSubmitAddress     string `json:"RestVoteSubmitAddress"`
	RestVoteSubmitPort        uint16 `json:"RestVoteSubmitPort"`
	RestVoteSubmitEndPoint    string `json:"RestVoteSubmitEndPoint"`
	RestVoteSubmitContentType string `json:"RestVoteSubmitContentType"`
}

type RequestBurstRangeClock struct {
	RestProposalRequestsRandomMinimumSeconds    uint8  `json:"RestProposalRequestsRandomMinimumSeconds"`
	RestProposalRequestsRandomMaximumSeconds    uint8  `json:"RestProposalRequestsRandomMaximumSeconds"`
	RestSubmitRequestsRandomMinimunMilliSeconds uint32 `json:"RestSubmitRequestsRandomMinimunMilliSeconds"`
	RestSubmitRequestsRandomMaximumMilliSeconds uint32 `json:"RestSubmitRequestsRandomMaximumMilliSeconds"`
}
