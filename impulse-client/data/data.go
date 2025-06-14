package data

import (
	"math/rand"
	"time"
)

type Vote struct {
	Topic           string `json:"topic"`
	DurationMinutes uint16 `json:"duration"`
}

type Topics struct {
	Votes []Vote `json:"topics"`
}

func (t *Topics) ShuffleTopics() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.Shuffle(len(t.Votes), func(i, j int) {
		t.Votes[i], t.Votes[j] = t.Votes[j], t.Votes[i]
	})
}

type BallotOption struct {
	BallotOptions []string `json:"BallotOptions"`
}

type UserData struct {
	UserHash string `json:"user_hash"`
}

type Users struct {
	UserLength int        `json:"length"`
	UserHashs  []UserData `json:"users"`
}

func (u *Users) GetUserHash(index int) string {
	if index < 0 || index >= len(u.UserHashs) {
		return ""
	}
	return u.UserHashs[index].UserHash
}

func (u *Users) ShuffleUserHashs() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.Shuffle(len(u.UserHashs), func(i, j int) {
		u.UserHashs[i], u.UserHashs[j] = u.UserHashs[j], u.UserHashs[i]
	})
}
