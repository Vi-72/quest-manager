package cmd

import (
	"net/http"

	httpmiddleware "quest-manager/internal/adapters/in/http/middleware"
)

// Middlewares returns the list of global HTTP middlewares configured for the API router.
func (cr *CompositionRoot) Middlewares() []func(http.Handler) http.Handler {
	middlewares := make([]func(http.Handler) http.Handler, 0, 4)

	if cr.authClient != nil {
		authMW := httpmiddleware.NewAuthMiddleware(cr.authClient)
		middlewares = append(middlewares, authMW.Auth)
	}

	return middlewares
}
