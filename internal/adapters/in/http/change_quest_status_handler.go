package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/quest"
)

// ChangeStatusRequest represents the JSON body for changing quest status.
type ChangeStatusRequest struct {
	Status quest.Status `json:"status"`
}

// ChangeQuestStatusHTTPHandler handles PATCH /api/v1/quests/{quest_id}/status requests.
type ChangeQuestStatusHTTPHandler struct {
	cmdHandler commands.ChangeQuestStatusCommandHandler
}

// NewChangeQuestStatusHTTPHandler creates a new instance of ChangeQuestStatusHTTPHandler.
func NewChangeQuestStatusHTTPHandler(handler commands.ChangeQuestStatusCommandHandler) *ChangeQuestStatusHTTPHandler {
	return &ChangeQuestStatusHTTPHandler{cmdHandler: handler}
}

// ServeHTTP processes the HTTP request to change the status of a quest.
func (h *ChangeQuestStatusHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var req ChangeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the status
	if !isValidStatus(req.Status) {
		http.Error(w, "invalid status value", http.StatusBadRequest)
		return
	}

	cmd := commands.ChangeQuestStatusCommand{
		ID:     questID,
		Status: req.Status,
	}

	result, err := h.cmdHandler.Handle(r.Context(), cmd)
	if err != nil {
		http.Error(w, "failed to change quest status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

// isValidStatus checks if the provided status is one of the allowed quest statuses.
func isValidStatus(status quest.Status) bool {
	switch status {
	case quest.StatusCreated,
		quest.StatusPosted,
		quest.StatusAssigned,
		quest.StatusInProgress,
		quest.StatusDeclined,
		quest.StatusCompleted:
		return true
	default:
		return false
	}
}
