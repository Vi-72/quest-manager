package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"

	"github.com/google/uuid"
)

// GetQuestByIDQueryHandler defines the interface for handling quest retrieval by ID.
type GetQuestByIDQueryHandler interface {
	Handle(ctx context.Context, questID uuid.UUID) (quest.Quest, error)
}

// getQuestByIDHandler is the implementation of GetQuestByIDQueryHandler.
type getQuestByIDHandler struct {
	questRepo ports.QuestRepository
}

// NewGetQuestByIDQueryHandler creates a new GetQuestByIDQueryHandler instance.
func NewGetQuestByIDQueryHandler(questRepo ports.QuestRepository) GetQuestByIDQueryHandler {
	return &getQuestByIDHandler{questRepo: questRepo}
}

// Handle processes the query to fetch a quest by its unique ID.
func (h *getQuestByIDHandler) Handle(ctx context.Context, questID uuid.UUID) (quest.Quest, error) {
	q, err := h.questRepo.GetByID(ctx, questID)
	if err != nil {
		// If quest not found, return NotFoundError for 404 response
		return quest.Quest{}, errs.NewNotFoundErrorWithCause("quest", questID.String(), err)
	}
	return q, nil
}
