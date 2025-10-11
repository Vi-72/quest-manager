package cmd

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

	// Middleware configuration
	Middleware MiddlewareConfig
}

// MiddlewareConfig contains configuration for HTTP middlewares
type MiddlewareConfig struct {
	DevAuth DevAuthConfig
}

// DevAuthConfig contains configuration for development/testing authentication
type DevAuthConfig struct {
	// Enabled enables development authentication mode (mock auth without gRPC)
	// If true: uses mock auth with static user ID from header
	// If false: uses real authentication via gRPC
	Enabled bool

	// HeaderName is the name of the HTTP header to read user ID from
	// Only used when Enabled=true
	// Default: "X-Dev-User-ID"
	HeaderName string

	// StaticUserID is the default user ID to use when header is not provided
	// Only used when Enabled=true
	// Default: "00000000-0000-0000-0000-000000000001"
	StaticUserID string
}
