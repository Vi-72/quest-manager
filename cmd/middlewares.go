package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"

	httpmiddleware "quest-manager/internal/adapters/in/http/middleware"
)

// Middlewares returns the list of global HTTP middlewares configured for the API router.
// Order matters: authentication first, then validation
func (c *Container) Middlewares(swagger *openapi3.T) []func(http.Handler) http.Handler {
	middlewares := make([]func(http.Handler) http.Handler, 0, 6)

	ctx := context.Background()

	// 1. Authentication middleware (first) - configurable
	if c.configs.Middleware.EnableAuth {
		if authClient := c.GetAuthClient(ctx); authClient != nil {
			authMW := httpmiddleware.NewAuthMiddleware(authClient)
			middlewares = append(middlewares, authMW.Auth)
			log.Printf("✅ Authentication middleware enabled")
		}
	} else {
		// Development mode - use mock authentication
		devHeader, devStatic := c.devAuthDefaults()
		middlewares = append(middlewares, c.mockAuthMiddleware(devHeader, devStatic))
		log.Printf("🔧 Mock authentication middleware enabled (dev mode)")
		log.Printf("   - Header name: %s", devHeader)
		log.Printf("   - Static user ID: %s", devStatic)
	}

	// 2. OpenAPI validation middleware (second) - always enabled
	if swagger != nil {
		validationMW, err := httpmiddleware.NewOpenAPIValidationMiddleware(swagger)
		if err == nil {
			middlewares = append(middlewares, validationMW.Validate)
		}
	}

	return middlewares
}

// mockAuthMiddleware returns a middleware for local development/testing
// that automatically authenticates requests with a configurable user ID
func (c *Container) mockAuthMiddleware(headerName, staticUserID string) func(http.Handler) http.Handler {
	// Pre-parse static user ID to validate it early
	defaultUserID, err := uuid.Parse(staticUserID)
	if err != nil {
		log.Fatalf("Invalid DEV_AUTH_STATIC_USER_ID: %s", staticUserID)
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to read user ID from header
			userIDStr := r.Header.Get(headerName)

			// If header is empty, use static user ID from config
			if userIDStr == "" {
				userIDStr = staticUserID
			}

			// Parse user ID
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				log.Printf("⚠️  Mock auth: invalid user ID '%s', using static: %s", userIDStr, staticUserID)
				userID = defaultUserID
			}

			// Add user ID to context
			ctx := httpmiddleware.UserIDToContext(r.Context(), userID)
			r = r.WithContext(ctx)

			log.Printf("🔧 Mock auth: request from user %s", userID)
			h.ServeHTTP(w, r)
		})
	}
}

// devAuthDefaults resolves development auth configuration with fallbacks to defaults.
func (c *Container) devAuthDefaults() (headerName, staticUserID string) {
	headerName = c.configs.Middleware.DevAuth.HeaderName
	if headerName == "" {
		headerName = DefaultDevAuthHeaderName
	}

	staticUserID = c.configs.Middleware.DevAuth.StaticUserID
	if staticUserID == "" {
		staticUserID = DefaultDevAuthStaticUserID
	}

	return headerName, staticUserID
}
