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
	uowFactory ports.UnitOfWorkFactory
}

// NewGetQuestByIDQueryHandler creates a new GetQuestByIDQueryHandler instance.
func NewGetQuestByIDQueryHandler(uowFactory ports.UnitOfWorkFactory) GetQuestByIDQueryHandler {
	return &getQuestByIDHandler{uowFactory: uowFactory}
}

// Handle processes the query to fetch a quest by its unique ID.
func (h *getQuestByIDHandler) Handle(ctx context.Context, questID uuid.UUID) (quest.Quest, error) {
	// Create a fresh UnitOfWork for this request
	uow, err := h.uowFactory.CreateUnitOfWork()
	if err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to create unit of work", err)
	}

	q, err := uow.QuestRepository().GetByID(ctx, questID)
	if err != nil {
		// If quest not found, return NotFoundError for 404 response
		return quest.Quest{}, errs.NewNotFoundErrorWithCause("quest", questID.String(), err)
	}
	return q, nil
}
