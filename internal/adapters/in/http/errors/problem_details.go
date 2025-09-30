package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ProblemDetails RFC 7807
type ProblemDetails struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (p *ProblemDetails) Error() string {
	return fmt.Sprintf("%d: %s - %s", p.Status, p.Title, p.Detail)
}

func (p *ProblemDetails) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		return
	}
}

// NewProblem creates a generic ProblemDetails with custom title and detail
func NewProblem(status int, title, detail string) *ProblemDetails {
	return &ProblemDetails{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
	}
}
