package cmd

import "quest-manager/internal/core/ports"

type Config struct {
	HttpPort            string
	DbHost              string
	DbPort              string
	DbUser              string
	DbPassword          string
	DbName              string
	DbSslMode           string
	EventGoroutineLimit int
	AuthGRPC            string

	// AuthClient is optional and used for testing to inject a mock auth client.
	// If provided, it will be used instead of creating a real gRPC auth client.
	AuthClient ports.AuthClient
}
