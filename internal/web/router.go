package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"quest-manager/cmd"
	"quest-manager/internal/generated/servers"
)

const apiV1Prefix = "/api/v1"

func NewRouter(root *cmd.CompositionRoot) http.Handler {
	router := chi.NewRouter()

	// Swagger JSON
	router.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		spec, err := servers.GetSwagger()
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
		w.Write(bytes)
	})

	// Swagger UI
	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
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
	apiHandler := servers.NewStrictHandler(strictHandler, nil)
	apiRouter := servers.HandlerFromMuxWithBaseURL(apiHandler, router, apiV1Prefix)

	apiRouter = servers.HandlerFromMux(apiHandler, router)
	return apiRouter
}
