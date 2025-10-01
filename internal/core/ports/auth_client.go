package ports

import (
	"context"

	"github.com/google/uuid"
)

type AuthClient interface {
	Authenticate(ctx context.Context, jwtToken string) (uuid.UUID, error)
}
