package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"quest-manager/cmd"
)

func main() {
	_ = godotenv.Load(".env")

	cfg := cmd.LoadConfig()

	router := cmd.NewRouter(cfg.JWTSecret)

	log.Printf("Server running on :%s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
