package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"
)

// ListQuestsQueryHandler defines the interface for handling quest listing.
// If status is nil, all quests are returned. Otherwise, filters by status.
type ListQuestsQueryHandler interface {
	Handle(ctx context.Context, status *quest.Status) ([]quest.Quest, error)
}

type listQuestsHandler struct {
	unitOfWorkFactory ports.UnitOfWorkFactory
}

// NewListQuestsQueryHandler creates a new ListQuestsQueryHandler instance.
func NewListQuestsQueryHandler(factory ports.UnitOfWorkFactory) ListQuestsQueryHandler {
	return &listQuestsHandler{unitOfWorkFactory: factory}
}

// Handle retrieves quests from the repository, optionally filtered by status.
func (h *listQuestsHandler) Handle(ctx context.Context, status *quest.Status) ([]quest.Quest, error) {
	unitOfWork, _, err := h.unitOfWorkFactory()
	if err != nil {
		return nil, errs.WrapInfrastructureError("failed to create unit of work for listing quests", err)
	}

	repo := unitOfWork.QuestRepository()
	if status != nil {
		// Validate status using domain logic - return DomainValidationError for 400
		if !quest.IsValidStatus(string(*status)) {
			return nil, errs.NewDomainValidationError("status", "must be one of 'created', 'posted', 'assigned', 'in_progress', 'declined', 'completed'")
		}

		// Filter by status
		return repo.FindByStatus(ctx, *status)
	}
	// Return all quests
	return repo.FindAll(ctx)
}
