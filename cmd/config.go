package cmd

import (
	"log"
	"os"
)

type Config struct {
	HTTPPort  string
	JWTSecret string
}

func LoadConfig() Config {
	return Config{
		HTTPPort:  getEnv("HTTP_PORT"),
		JWTSecret: getEnv("JWT_SECRET"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing environment variable %s", key)
	}
	return val
}
