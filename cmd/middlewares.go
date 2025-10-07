package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

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
		log.Printf("ℹ️  Authentication middleware disabled by configuration")
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
