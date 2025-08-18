package cmd

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func NewRouter(jwtSecret string) http.Handler {
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	r := chi.NewRouter()

	r.Post("/auth/login", loginHandler(tokenAuth))
	r.Post("/auth/register", registerHandler)

	// Example of protected routes using JWT middleware
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		// r.Get("/protected", protectedHandler) // placeholder
	})

	return r
}

func loginHandler(tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, token, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 1})
		w.Write([]byte(token))
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("registered"))
}
