package cmd

import (
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

	// 1. Authentication middleware (first) - always enabled
	// Mode: Dev (mock) or Production (real gRPC)
	if c.configs.Middleware.DevAuth.Enabled {
		// Development mode - use mock authentication
		devHeader, devStatic := c.devAuthDefaults()
		middlewares = append(middlewares, c.mockAuthMiddleware(devHeader, devStatic))
		log.Printf("🔧 Development authentication mode enabled")
		log.Printf("   - Header name: %s", devHeader)
		log.Printf("   - Static user ID: %s", devStatic)
	} else {
		// Production mode - use real auth client
		if authClient := c.GetAuthClient(); authClient != nil {
			authMW := httpmiddleware.NewAuthMiddleware(authClient)
			middlewares = append(middlewares, authMW.Auth)
			log.Printf("Production authentication enabled (gRPC)")
		} else {
			log.Fatal("Production auth enabled but failed to create auth client")
		}
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
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to read user ID from header
			userIDStr := r.Header.Get(headerName)

			// If header is empty, use static user ID from config
			if userIDStr == "" {
				userIDStr = staticUserID
			}

			// Parse user ID
			userID, _ := uuid.Parse(userIDStr)

			// Add user ID to context
			ctx := httpmiddleware.UserIDToContext(r.Context(), userID)
			r = r.WithContext(ctx)

			log.Printf("Mock auth: request from user %s", userID)
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
