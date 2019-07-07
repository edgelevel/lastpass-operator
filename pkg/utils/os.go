package utils

import (
	"fmt"
	"log"
	"os"
)

func GetEnvOrDie(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.Printf("Environment variable found: [%s=%s]", key, value)
		return value
	} else {
		panic(fmt.Sprintf("No environment variable found with name: %s", key))
	}
}