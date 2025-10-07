package cmd

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	httpmiddleware "quest-manager/internal/adapters/in/http/middleware"
)

// Middlewares returns the list of global HTTP middlewares configured for the API router.
// Order matters: authentication first, then validation
func (cr *CompositionRoot) Middlewares(swagger *openapi3.T) []func(http.Handler) http.Handler {
	middlewares := make([]func(http.Handler) http.Handler, 0, 4)

	// 1. Authentication middleware (first)
	if cr.authClient != nil {
		authMW := httpmiddleware.NewAuthMiddleware(cr.authClient)
		middlewares = append(middlewares, authMW.Auth)
	}

	// 2. OpenAPIs validation middleware (second)
	if swagger != nil {
		validationMW, err := httpmiddleware.NewOpenAPIValidationMiddleware(swagger)
		if err == nil {
			middlewares = append(middlewares, validationMW.Validate)
		}
	}

	return middlewares
}
