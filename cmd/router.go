package cmd

import (
	"errors"
	"net/http"

	"quest-manager/internal/adapters/in/http/validations"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware" // ← добавлено

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http/problems"
	"quest-manager/internal/pkg/errs"
)

const apiV1Prefix = "/api/v1"

func NewRouter(root *CompositionRoot) http.Handler {
	router := chi.NewRouter()

	// --- Базовые middleware ---
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// --- Health check ---
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Swagger JSON
	router.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		spec, err := v1.GetSwagger()
		if err != nil {
			http.Error(w, "failed to load OpenAPI spec: "+err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := spec.MarshalJSON()
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
			problem := problems.NewBadRequest("Invalid request parameters: " + err.Error())
			problem.WriteResponse(w)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			// Check if it's a validation error from our pkg/validations
			var validationErr *validations.ValidationError
			if errors.As(err, &validationErr) {
				// Convert validation error to Problem Details and return
				problem := validations.ConvertValidationErrorToProblem(validationErr)
				problem.WriteResponse(w)
				return
			}

			// Check if it's a domain validation error from application layer
			var domainValidationErr *errs.DomainValidationError
			if errors.As(err, &domainValidationErr) {
				// Convert to 400 Bad Request
				problem := validations.ConvertDomainValidationErrorToProblem(domainValidationErr)
				problem.WriteResponse(w)
				return
			}

			// Check if it's a not found error from application layer
			var notFoundErr *errs.NotFoundError
			if errors.As(err, &notFoundErr) {
				// Convert to 404 Not Found
				problem := validations.ConvertNotFoundErrorToProblem(notFoundErr)
				problem.WriteResponse(w)
				return
			}

			// Handle other response errors
			problem := problems.NewBadRequest("Response error: " + err.Error())
			problem.WriteResponse(w)
		},
	})

	apiRouter := v1.HandlerFromMuxWithBaseURL(apiHandler, router, apiV1Prefix)

	return apiRouter
}
