package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"quest-manager/internal/core/application/usecases/queries"
)

// GetQuestByIDHTTPHandler handles GET /api/v1/quests/{quest_id} requests.
type GetQuestByIDHTTPHandler struct {
	queryHandler queries.GetQuestByIDQueryHandler
}

// NewGetQuestByIDHTTPHandler creates a new instance of GetQuestByIDHTTPHandler.
func NewGetQuestByIDHTTPHandler(handler queries.GetQuestByIDQueryHandler) *GetQuestByIDHTTPHandler {
	return &GetQuestByIDHTTPHandler{queryHandler: handler}
}

// ServeHTTP processes the HTTP request to get a quest by its ID.
func (h *GetQuestByIDHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questIDStr, ok := vars["quest_id"]
	if !ok {
		http.Error(w, "quest_id path parameter is required", http.StatusBadRequest)
		return
	}

	questID, err := uuid.Parse(questIDStr)
	if err != nil {
		http.Error(w, "invalid quest_id format", http.StatusBadRequest)
		return
	}

	query := queries.GetQuestByIDQuery{ID: questID}
	result, err := h.queryHandler.Handle(r.Context(), query)
	if err != nil {
		http.Error(w, "quest not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result.Quest)
}
