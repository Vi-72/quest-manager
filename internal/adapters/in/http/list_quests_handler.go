package http

import (
	"encoding/json"
	"net/http"

	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/quest"
)

// ListQuestsHTTPHandler handles GET /api/v1/quests requests.
type ListQuestsHTTPHandler struct {
	queryHandler queries.ListQuestsQueryHandler
}

// NewListQuestsHTTPHandler creates a new instance of ListQuestsHTTPHandler.
func NewListQuestsHTTPHandler(handler queries.ListQuestsQueryHandler) *ListQuestsHTTPHandler {
	return &ListQuestsHTTPHandler{queryHandler: handler}
}

// ServeHTTP processes the HTTP request to list quests with an optional status filter.
func (h *ListQuestsHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statusParam := r.URL.Query().Get("status")

	var listQuery queries.ListQuestsQuery
	if statusParam != "" {
		status := quest.Status(statusParam)
		// Validate status
		if !isValidStatus(status) {
			http.Error(w, "invalid status value", http.StatusBadRequest)
			return
		}
		listQuery.Status = &status
	}

	result, err := h.queryHandler.Handle(r.Context(), listQuery)
	if err != nil {
		http.Error(w, "failed to list quests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result.Quests)
}
