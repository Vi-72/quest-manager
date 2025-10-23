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
	uowFactory ports.UnitOfWorkFactory
}

// NewListAssignedQuestsQueryHandler creates a new instance of ListAssignedQuestsQueryHandler.
func NewListAssignedQuestsQueryHandler(uowFactory ports.UnitOfWorkFactory) ListAssignedQuestsQueryHandler {
	return &listAssignedQuestsHandler{uowFactory: uowFactory}
}

// Handle retrieves all quests assigned to the given user.
func (h *listAssignedQuestsHandler) Handle(ctx context.Context, userID uuid.UUID) ([]quest.Quest, error) {
	uow, err := h.uowFactory.CreateUnitOfWork()
	if err != nil {
		return nil, err
	}
	return uow.QuestRepository().FindByAssignee(ctx, userID)
}
