package utils

import (
	"log"
	"os"
)

func MustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env var %s", key)
	}
	return v
}
