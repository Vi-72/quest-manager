package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	httperrors "quest-manager/internal/adapters/in/http/errors"
	"quest-manager/internal/adapters/out/client/auth"

	"github.com/google/uuid"
)

type contextKey string

func (c contextKey) String() string {
	return "context key " + string(c)
}

var (
	contextKeyAuthenticatedUser = contextKey("authenticated_user")
)

// UserIDToContext adds user ID to context (exported for testing/mocking)
func UserIDToContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKeyAuthenticatedUser, userID)
}

// UserIDFromContext retrieves user ID from context
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	c, ok := ctx.Value(contextKeyAuthenticatedUser).(uuid.UUID)
	return c, ok
}

type AuthMiddleware struct {
	authClient auth.Client
}

func NewAuthMiddleware(authClient auth.Client) *AuthMiddleware {
	return &AuthMiddleware{authClient: authClient}
}

func (mw *AuthMiddleware) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		jwtStr, err := bearerTokenFromHeader(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := mw.authClient.Authenticate(ctx, jwtStr)
		if err != nil {
			// Handle token expired separately for better error messages
			if errors.Is(err, auth.ErrTokenExpired) {
				problem := httperrors.NewProblem(
					http.StatusUnauthorized,
					"Token Expired",
					"JWT token has expired, please refresh your token",
				)
				problem.WriteResponse(w)
				return
			}

			// Generic authentication failure
			problem := httperrors.NewProblem(
				http.StatusUnauthorized,
				"Authentication Failed",
				"Invalid or malformed authentication token",
			)
			problem.WriteResponse(w)
			return
		}

		r = r.WithContext(UserIDToContext(ctx, userID))
		h.ServeHTTP(w, r)
	})
}

// bearerTokenFromHeader extracts the JWT token from the Authorization header.
func bearerTokenFromHeader(r *http.Request) (string, error) {
	const bearerPrefix = "Bearer "
	authHeader := r.Header.Get("Authorization")

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("missing or invalid Authorization header")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if token == "" {
		return "", errors.New("missing or invalid Authorization header")
	}

	return token, nil
}
