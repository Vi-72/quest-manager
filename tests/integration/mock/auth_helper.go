package mock

import (
	"context"

	"github.com/google/uuid"
)

// AlwaysSuccessAuthClient is a simple wrapper that always returns successful authentication.
// This is useful for integration tests where we don't need gomock.Controller complexity.
type AlwaysSuccessAuthClient struct {
	DefaultUserID uuid.UUID
}

// NewAlwaysSuccessAuthClient creates a new auth client that always succeeds with default user ID.
func NewAlwaysSuccessAuthClient() *AlwaysSuccessAuthClient {
	return &AlwaysSuccessAuthClient{
		DefaultUserID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	}
}

// Authenticate always returns the configured user ID without any validation.
func (a *AlwaysSuccessAuthClient) Authenticate(ctx context.Context, jwtToken string) (uuid.UUID, error) {
	return a.DefaultUserID, nil
}
