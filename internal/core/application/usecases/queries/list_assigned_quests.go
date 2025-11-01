package queries

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// ListAssignedQuestsQueryHandler defines the interface for handling assigned quests retrieval.
type ListAssignedQuestsQueryHandler interface {
	Handle(ctx context.Context, userID uuid.UUID) ([]quest.Quest, error)
}

type listAssignedQuestsHandler struct {
	questRepo ports.QuestRepository
}

// NewListAssignedQuestsQueryHandler creates a new instance of ListAssignedQuestsQueryHandler.
func NewListAssignedQuestsQueryHandler(questRepo ports.QuestRepository) ListAssignedQuestsQueryHandler {
	return &listAssignedQuestsHandler{questRepo: questRepo}
}

// Handle retrieves all quests assigned to the given user.
func (h *listAssignedQuestsHandler) Handle(ctx context.Context, userID uuid.UUID) ([]quest.Quest, error) {
	return h.questRepo.FindByAssignee(ctx, userID)
}
