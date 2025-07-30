package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"

	"github.com/google/uuid"
)

// GetQuestByIDQueryHandler defines the interface for handling quest retrieval by ID.
type GetQuestByIDQueryHandler interface {
	Handle(ctx context.Context, questID uuid.UUID) (quest.Quest, error)
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
func (h *getQuestByIDHandler) Handle(ctx context.Context, questID uuid.UUID) (quest.Quest, error) {
	return h.repo.GetByID(ctx, questID)
}
