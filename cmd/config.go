package cmd

import "quest-manager/internal/adapters/out/client/auth"

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

	// AuthFactory is used to create auth clients
	AuthFactory *auth.Factory

	// Middleware configuration
	Middleware MiddlewareConfig
}

// MiddlewareConfig contains configuration for HTTP middlewares
type MiddlewareConfig struct {
	// EnableAuth enables authentication middleware
	EnableAuth bool
}
