package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("ERROR: Could not load .env file. Proceeding without it (assuming environment variables are set externally): %v", err)
		log.Fatalf("CRITICAL ERROR: Failed to load .env file: %v", err)
	} else {
		log.Println("INFO: Environment variables successfully loaded from .env file.")
	}

	serverName := GetEnvVar("SERVER_NAME")

	log.Println("=======================================================================================")
	log.Printf("                  Welcome! Starting %s", serverName)
	log.Println("=======================================================================================")
}

func GetEnvVar(varName string) string {
	value := os.Getenv(varName)

	if value == "" {
		log.Fatalf("CRITICAL ERROR: Required environment variable \"%s\" is not defined. Application will exit.", varName)
	}

	return value
}

func GetIntEnvVar(varName string) int {
	stringValue := GetEnvVar(varName)
	intValue, err := strconv.Atoi(stringValue)

	if err != nil {
		log.Fatalf("CRITICAL ERROR: Environment variable \"%s\" (\"%s\") is not a valid integer. Application will exit. Error: %v", varName, stringValue, err)
	}

	return intValue
}
