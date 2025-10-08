package cmd

import (
	"quest-manager/internal/adapters/out/client/auth"
)

const (
	// DefaultDevAuthHeaderName is the default header name for dev auth
	DefaultDevAuthHeaderName = "X-Dev-User-ID"

	// DefaultDevAuthStaticUserID is the default static user ID for dev auth
	DefaultDevAuthStaticUserID = "00000000-0000-0000-0000-000000000001"
)

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

	// DevAuth contains configuration for development authentication (when EnableAuth=false)
	DevAuth DevAuthConfig
}

// DevAuthConfig contains configuration for development/testing authentication
type DevAuthConfig struct {
	// HeaderName is the name of the HTTP header to read user ID from
	// Default: "X-Dev-User-ID"
	HeaderName string

	// StaticUserID is the default user ID to use when header is not provided
	// Default: "00000000-0000-0000-0000-000000000001"
	StaticUserID string
}
