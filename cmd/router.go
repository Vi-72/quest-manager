package cmd

import (
	"errors"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	v1 "quest-manager/api/http/quests/v1"
	httperrors "quest-manager/internal/adapters/in/http/errors"
	"quest-manager/internal/pkg/errs"
)

const apiV1Prefix = "/api/v1"

func NewRouter(root *CompositionRoot) http.Handler {
	router := chi.NewRouter()

	// --- Базовые middleware ---
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Logger)

	swagger, err := v1.GetSwagger()
	if err != nil {
		panic("failed to load OpenAPI spec: " + err.Error())
	}

	swagger.Servers = []*openapi3.Server{{URL: apiV1Prefix}}

	// --- Health check ---
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Swagger JSON
	router.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := swagger.MarshalJSON()
		if err != nil {
			http.Error(w, "failed to marshal OpenAPI: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bytes)
	})

	// Swagger UI
	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
			  <meta charset="UTF-8">
			  <title>Swagger UI</title>
			  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
			</head>
			<body>
			  <div id="swagger-ui"></div>
			  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
			  <script>
				window.onload = () => {
				  SwaggerUIBundle({
					url: "/openapi.json",
					dom_id: "#swagger-ui",
				  });
				};
			  </script>
			</body>
			</html>
		`))
	})

	strictHandler := root.NewApiHandler()

	// Create StrictHandler with custom error handling for parameter parsing and validation
	apiHandler := v1.NewStrictHandlerWithOptions(strictHandler, nil, v1.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			// Handle parameter parsing errors with detailed messages
			problem := httperrors.NewBadRequest("Invalid request parameters: " + err.Error())
			problem.WriteResponse(w)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			// Check if it's a domain validation error from application layer
			var domainValidationErr *errs.DomainValidationError
			if errors.As(err, &domainValidationErr) {
				// Convert to 400 Bad Request
				problem := httperrors.NewDomainValidationProblem(domainValidationErr)
				problem.WriteResponse(w)
				return
			}

			// Check if it's a not found error from application layer
			var notFoundErr *errs.NotFoundError
			if errors.As(err, &notFoundErr) {
				// Convert to 404 Not Found
				problem := httperrors.NewNotFoundProblem(notFoundErr)
				problem.WriteResponse(w)
				return
			}

			// Handle other response errors
			problem := httperrors.NewBadRequest("Response error: " + err.Error())
			problem.WriteResponse(w)
		},
	})

	apiRouter := chi.NewRouter()

	// Apply middlewares in order: authentication first, then validation
	for _, mw := range root.Middlewares(swagger) {
		apiRouter.Use(mw)
	}

	v1.HandlerFromMux(apiHandler, apiRouter)

	router.Mount(apiV1Prefix, apiRouter)

	return router
}
