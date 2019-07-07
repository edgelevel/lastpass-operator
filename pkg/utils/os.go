package utils

import (
	"fmt"
	"os"
)

func GetEnvOrDie(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		panic(fmt.Sprintf("No environment variable found with name: %s", key))
	}
}
