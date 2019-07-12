package utils

import (
	"fmt"
	"log"
	"os"
)

// GetEnvOrDie retrieve an environment variable or exit
func GetEnvOrDie(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.Printf("Found environment variable: [%s]", key)
		return value
	}
	panic(fmt.Sprintf("No environment variable found: [%s]", key))
}
