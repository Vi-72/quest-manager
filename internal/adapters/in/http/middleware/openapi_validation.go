package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"

	httperrors "quest-manager/internal/adapters/in/http/errors"
)

// OpenAPIValidationMiddleware validates HTTP requests against OpenAPI specification
type OpenAPIValidationMiddleware struct {
	router routers.Router
}

// NewOpenAPIValidationMiddleware creates a new OpenAPI validation middleware
func NewOpenAPIValidationMiddleware(doc *openapi3.T) (*OpenAPIValidationMiddleware, error) {
	router, err := legacy.NewRouter(doc)
	if err != nil {
		return nil, err
	}

	return &OpenAPIValidationMiddleware{
		router: router,
	}, nil
}

// Validate validates HTTP requests against OpenAPI specification
func (mw *OpenAPIValidationMiddleware) Validate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := mw.validateRequest(r); err != nil {
			problem := httperrors.NewBadRequest("Request validation failed: " + err.Error())
			problem.WriteResponse(w)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// validateRequest validates the HTTP request against OpenAPI specification
func (mw *OpenAPIValidationMiddleware) validateRequest(r *http.Request) error {
	var (
		bodyBytes []byte
		err       error
	)

	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		if err := r.Body.Close(); err != nil {
			return err
		}
	}

	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	requestForValidation := r.Clone(r.Context())
	if requestForValidation.URL != nil {
		clonedURL := *requestForValidation.URL
		requestForValidation.URL = &clonedURL
	}
	if requestForValidation.URL != nil {
		requestForValidation.URL.RawPath = requestForValidation.URL.Path
	}
	requestForValidation.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	route, pathParams, err := mw.router.FindRoute(requestForValidation)
	if err != nil {
		return err
	}

	validationInput := &openapi3filter.RequestValidationInput{
		Request:    requestForValidation,
		PathParams: pathParams,
		Route:      route,
		Options: &openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
			MultiError:         true,
		},
	}

	if err := openapi3filter.ValidateRequest(r.Context(), validationInput); err != nil {
		return err
	}

	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return nil
}
