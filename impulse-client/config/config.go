package config

type Vote struct {
	Topic           string `json:"topic"`
	DurationMinutes uint16 `json:"duration"`
}

type Topics struct {
	Votes []Vote `json:"topics"`
}

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
