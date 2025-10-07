package cmd

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/Vi-72/quest-auth/api/grpc/sdk/go/auth/v1"

	authclient "quest-manager/internal/adapters/out/client/auth"
	"quest-manager/internal/core/ports"
)

// AuthClientFactory creates and configures authentication clients
type AuthClientFactory struct {
	configs Config
}

// NewAuthClientFactory creates a new AuthClientFactory
func NewAuthClientFactory(configs Config) *AuthClientFactory {
	return &AuthClientFactory{
		configs: configs,
	}
}

// CreateAuthClient creates an authentication client based on configuration
// Returns the client and a closer function for cleanup
func (f *AuthClientFactory) CreateAuthClient() (ports.AuthClient, Closer) {
	// If AuthClient is provided in config (e.g., mock for tests), use it
	if f.configs.AuthClient != nil {
		return f.configs.AuthClient, nil
	}

	// Otherwise, wire Auth gRPC client (optional: if AUTH_GRPC provided)
	if addr := f.configs.AuthGRPC; addr != "" {
		// Use grpc.NewClient instead of deprecated DialContext
		// NewClient is lazy - it connects on first RPC call
		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to create auth gRPC client at %s: %v", addr, err)
		}

		authSDKClient := authv1.NewAuthServiceClient(conn)
		authClient := authclient.NewUserAuthClient(authSDKClient)

		return authClient, connCloser{conn}
	}

	// No auth client configured
	return nil, nil
}

// connCloser implements Closer interface for gRPC connections
type connCloser struct{ *grpc.ClientConn }

func (c connCloser) Close() error { return c.ClientConn.Close() }
