package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(file string) error {
	if err := godotenv.Load(file); err != nil {
		return fmt.Errorf("error loading env file %s: %w", file, err)
	}
	return nil
}

func QdrantHost() string {
	host := os.Getenv("QDRANT_HOST")
	if host == "" {
		host = "http://localhost:6333"
	}
	return host
}
