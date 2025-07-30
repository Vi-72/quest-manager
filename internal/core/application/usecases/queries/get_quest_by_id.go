package queries

import (
	"context"

	"github.com/google/uuid"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// GetQuestByIDQuery represents the input data required to fetch a quest by its ID.
type GetQuestByIDQuery struct {
	ID uuid.UUID
}

// GetQuestByIDResult represents the result of fetching a quest by its ID.
type GetQuestByIDResult struct {
	Quest quest.Quest
}

// GetQuestByIDQueryHandler defines the interface for handling GetQuestByIDQuery.
type GetQuestByIDQueryHandler interface {
	Handle(ctx context.Context, query GetQuestByIDQuery) (GetQuestByIDResult, error)
}

// getQuestByIDHandler is the implementation of GetQuestByIDQueryHandler.
type getQuestByIDHandler struct {
	repo ports.QuestRepository
}

// NewGetQuestByIDQueryHandler creates a new GetQuestByIDQueryHandler instance.
func NewGetQuestByIDQueryHandler(repo ports.QuestRepository) GetQuestByIDQueryHandler {
	return &getQuestByIDHandler{repo: repo}
}

// Handle processes the query to fetch a quest by its unique ID.
func (h *getQuestByIDHandler) Handle(ctx context.Context, query GetQuestByIDQuery) (GetQuestByIDResult, error) {
	q, err := h.repo.GetByID(ctx, query.ID)
	if err != nil {
		return GetQuestByIDResult{}, err
	}

	return GetQuestByIDResult{Quest: q}, nil
}
