//go:generate mockgen -destination=../mock/auth.go -package=mock quest-manager/internal/adapters/out/client/auth Client
package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"

	authv1 "github.com/Vi-72/quest-auth/api/grpc/sdk/go/auth/v1"
)

var _ Client = &client{}

var (
	// ErrTokenExpired is returned when the auth service reports an expired JWT.
	ErrTokenExpired = errors.New("token expired")
	// errMissingUser signals that the auth service response does not contain user data.
	errMissingUser = errors.New("auth response missing user")
)

type Client interface {
	Authenticate(ctx context.Context, jwtToken string) (uuid.UUID, error)
}

type client struct {
	authClient authv1.AuthServiceClient
}

func NewUserAuthClient(authClient authv1.AuthServiceClient) Client {
	return &client{authClient: authClient}
}

func (c *client) Authenticate(ctx context.Context, jwtToken string) (uuid.UUID, error) {
	if strings.TrimSpace(jwtToken) == "" {
		return uuid.Nil, errors.New("jwt token is empty")
	}

	resp, err := c.authClient.Authenticate(ctx, &authv1.AuthenticateRequest{JwtToken: jwtToken})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			tokenExpiredMsg := "token is expired"
			if strings.Contains(strings.ToLower(s.Message()), tokenExpiredMsg) {
				return uuid.Nil, ErrTokenExpired
			}
		}
		return uuid.Nil, fmt.Errorf("authenticate grpc call: %w", err)
	}

	user := resp.GetUser()
	if user == nil {
		return uuid.Nil, errMissingUser
	}

	userID, err := uuid.Parse(user.GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
	}

	return userID, nil
}
