package util

import (
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/andantan/vote-blockchain-server/impulse-client/data"
)

func RandRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandOption(params []string) string {
	return string(params[rand.Intn(len(params))])
}

func RandomHashString() string {
	return hex.EncodeToString(ToSlice(HashFromBytes(RandomBytes(32))))
}

func RandomBytes(size int) []byte {
	ticket := make([]byte, size)

	crand.Read(ticket)

	return ticket
}

func HashFromBytes(b []byte) [32]uint8 {
	var t [32]uint8

	for i := range 32 {
		t[i] = b[i]
	}

	return t
}

func ToSlice(d [32]uint8) []byte {
	b := make([]byte, 32)

	for i := range 32 {
		b[i] = d[i]
	}

	return b
}

func MakeRandomUsers(numUser int) {
	userHashs := make([]data.UserData, 0, numUser)

	for range numUser {
		userData := data.UserData{
			UserHash: RandomHashString(),
		}

		userHashs = append(userHashs, userData)
	}

	usersFileContent := data.Users{
		UserLength: numUser,
		UserHashs:  userHashs,
	}

	jsonData, err := json.MarshalIndent(usersFileContent, "", "    ")

	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	path := filepath.Join("./", "data", "user_data.json")

	file, err := os.Create(path)

	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}

	defer file.Close()

	_, err = file.Write(jsonData)

	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}
}
