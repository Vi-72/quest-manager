package queries

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"
)

// ListAssignedQuestsQueryHandler defines the interface for handling assigned quests retrieval.
type ListAssignedQuestsQueryHandler interface {
	Handle(ctx context.Context, userID uuid.UUID) ([]quest.Quest, error)
}

type listAssignedQuestsHandler struct {
	unitOfWorkFactory ports.UnitOfWorkFactory
}

// NewListAssignedQuestsQueryHandler creates a new instance of ListAssignedQuestsQueryHandler.
func NewListAssignedQuestsQueryHandler(factory ports.UnitOfWorkFactory) ListAssignedQuestsQueryHandler {
	return &listAssignedQuestsHandler{unitOfWorkFactory: factory}
}

// Handle retrieves all quests assigned to the given user.
func (h *listAssignedQuestsHandler) Handle(ctx context.Context, userID uuid.UUID) ([]quest.Quest, error) {
	unitOfWork, _, err := h.unitOfWorkFactory()
	if err != nil {
		return nil, errs.WrapInfrastructureError("failed to create unit of work for assigned quests lookup", err)
	}

	return unitOfWork.QuestRepository().FindByAssignee(ctx, userID)
}
