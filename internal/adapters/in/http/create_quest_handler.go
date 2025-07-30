package http

import (
	"encoding/json"
	"net/http"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

// CreateQuestRequest represents the JSON payload for creating a quest.
type CreateQuestRequest struct {
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Difficulty         string   `json:"difficulty"`
	Reward             string   `json:"reward"`
	TargetLatitude     float64  `json:"target_latitude"`
	TargetLongitude    float64  `json:"target_longitude"`
	ExecutionLatitude  float64  `json:"execution_latitude"`
	ExecutionLongitude float64  `json:"execution_longitude"`
	Equipment          []string `json:"equipment"`
	Skills             []string `json:"skills"`
	Creator            string   `json:"creator"`
}

// CreateQuestHTTPHandler handles the POST /api/v1/quests endpoint.
type CreateQuestHTTPHandler struct {
	cmdHandler commands.CreateQuestCommandHandler
}

// NewCreateQuestHTTPHandler creates a new instance of CreateQuestHTTPHandler.
func NewCreateQuestHTTPHandler(handler commands.CreateQuestCommandHandler) *CreateQuestHTTPHandler {
	return &CreateQuestHTTPHandler{cmdHandler: handler}
}

// ServeHTTP processes the HTTP request for creating a quest.
func (h *CreateQuestHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req CreateQuestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	targetLocation, err := kernel.NewGeoCoordinate(req.TargetLatitude, req.TargetLongitude)
	if err != nil {
		http.Error(w, "invalid target location: "+err.Error(), http.StatusBadRequest)
		return
	}

	executionLocation, err := kernel.NewGeoCoordinate(req.ExecutionLatitude, req.ExecutionLongitude)
	if err != nil {
		http.Error(w, "invalid execution location: "+err.Error(), http.StatusBadRequest)
		return
	}

	difficulty := quest.Difficulty(req.Difficulty)
	if difficulty != quest.DifficultyEasy && difficulty != quest.DifficultyMedium && difficulty != quest.DifficultyHard {
		http.Error(w, "invalid difficulty level", http.StatusBadRequest)
		return
	}

	cmd := commands.CreateQuestCommand{
		Title:             req.Title,
		Description:       req.Description,
		Difficulty:        difficulty,
		Reward:            req.Reward,
		TargetLocation:    targetLocation,
		ExecutionLocation: executionLocation,
		Equipment:         req.Equipment,
		Skills:            req.Skills,
		Creator:           req.Creator,
	}

	result, err := h.cmdHandler.Handle(r.Context(), cmd)
	if err != nil {
		http.Error(w, "failed to create quest: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}
