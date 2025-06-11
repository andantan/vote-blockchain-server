package data

type Vote struct {
	Topic           string `json:"topic"`
	DurationMinutes uint16 `json:"duration"`
}

type Topics struct {
	Votes []Vote `json:"topics"`
}

type BallotOption struct {
	BallotOptions []string `json:"BallotOptions"`
}
