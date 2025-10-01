package commands

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"
)

// ChangeQuestStatusCommandHandler defines the interface for handling quest status changes.
type ChangeQuestStatusCommandHandler interface {
	Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (ChangeQuestStatusResult, error)
}

type changeQuestStatusHandler struct {
	unitOfWorkFactory ports.UnitOfWorkFactory
}

// NewChangeQuestStatusCommandHandler creates a new ChangeQuestStatusCommandHandler instance.
func NewChangeQuestStatusCommandHandler(factory ports.UnitOfWorkFactory) ChangeQuestStatusCommandHandler {
	return &changeQuestStatusHandler{
		unitOfWorkFactory: factory,
	}
}

// Handle updates the quest status with validation and domain business rules.
func (h *changeQuestStatusHandler) Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (ChangeQuestStatusResult, error) {
	// Validate status - this is domain validation error → 400
	if !quest.IsValidStatus(string(cmd.Status)) {
		return ChangeQuestStatusResult{}, errs.NewDomainValidationError("status", "must be one of 'created', 'posted', 'assigned', 'in_progress', 'declined', 'completed'")
	}

	unitOfWork, eventPublisher, err := h.unitOfWorkFactory()
	if err != nil {
		return ChangeQuestStatusResult{}, errs.WrapInfrastructureError("failed to create quest status unit of work", err)
	}

	var (
		txStarted bool
		committed bool
	)
	defer func() {
		if txStarted && !committed {
			_ = unitOfWork.Rollback()
		}
	}()

	// Begin transaction
	if err := unitOfWork.Begin(ctx); err != nil {
		return ChangeQuestStatusResult{}, errs.WrapInfrastructureError("failed to begin quest status change transaction", err)
	}
	txStarted = true

	// Get quest - if not found → 404
	q, err := unitOfWork.QuestRepository().GetByID(ctx, cmd.QuestID)
	if err != nil {
		return ChangeQuestStatusResult{}, errs.NewNotFoundErrorWithCause("quest", cmd.QuestID.String(), err)
	}

	// Use domain logic for status change - domain validation error → 400
	if err := q.ChangeStatus(cmd.Status); err != nil {
		return ChangeQuestStatusResult{}, errs.NewDomainValidationErrorWithCause("status", "invalid status transition", err)
	}

	// Save quest - infrastructure error → 500
	if err := unitOfWork.QuestRepository().Save(ctx, q); err != nil {
		return ChangeQuestStatusResult{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Publish domain events within the same transaction
	if eventPublisher != nil {
		if err := eventPublisher.Publish(ctx, q.GetDomainEvents()...); err != nil {
			return ChangeQuestStatusResult{}, errs.WrapInfrastructureError("failed to publish events", err)
		}
	}

	// Commit transaction
	if err := unitOfWork.Commit(ctx); err != nil {
		return ChangeQuestStatusResult{}, errs.WrapInfrastructureError("failed to commit quest status change transaction", err)
	}
	committed = true

	// Clear events after successful commit
	q.ClearDomainEvents()

	// Form result from updated quest
	return ChangeQuestStatusResult{
		ID:       q.ID(),
		Assignee: q.Assignee, // Now both are *uuid.UUID
		Status:   string(q.Status),
	}, nil
}
