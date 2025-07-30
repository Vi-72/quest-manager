package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"quest-manager/internal/core/application/usecases/commands"
)

// AssignQuestRequest represents the JSON payload for assigning a quest.
type AssignQuestRequest struct {
	UserID string `json:"user_id"`
}

// AssignQuestHTTPHandler handles POST /api/v1/quests/{quest_id}/assign requests.
type AssignQuestHTTPHandler struct {
	cmdHandler commands.AssignQuestCommandHandler
}

// NewAssignQuestHTTPHandler creates a new instance of AssignQuestHTTPHandler.
func NewAssignQuestHTTPHandler(handler commands.AssignQuestCommandHandler) *AssignQuestHTTPHandler {
	return &AssignQuestHTTPHandler{cmdHandler: handler}
}

// ServeHTTP processes the HTTP request to assign a quest to a user.
func (h *AssignQuestHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var req AssignQuestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	cmd := commands.AssignQuestCommand{
		ID:     questID,
		UserID: req.UserID,
	}

	result, err := h.cmdHandler.Handle(r.Context(), cmd)
	if err != nil {
		http.Error(w, "failed to assign quest: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
