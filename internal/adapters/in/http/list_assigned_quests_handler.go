package http

import (
	"encoding/json"
	"net/http"

	"quest-manager/internal/core/application/usecases/queries"
)

// ListAssignedQuestsHTTPHandler handles GET /api/v1/quests/assigned requests.
type ListAssignedQuestsHTTPHandler struct {
	queryHandler queries.ListAssignedQuestsQueryHandler
}

// NewListAssignedQuestsHTTPHandler creates a new instance of ListAssignedQuestsHTTPHandler.
func NewListAssignedQuestsHTTPHandler(handler queries.ListAssignedQuestsQueryHandler) *ListAssignedQuestsHTTPHandler {
	return &ListAssignedQuestsHTTPHandler{queryHandler: handler}
}

// ServeHTTP processes the HTTP request to get quests assigned to a specific user.
func (h *ListAssignedQuestsHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	query := queries.ListAssignedQuestsQuery{UserID: userID}
	result, err := h.queryHandler.Handle(r.Context(), query)
	if err != nil {
		http.Error(w, "failed to get assigned quests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result.Quests)
}
