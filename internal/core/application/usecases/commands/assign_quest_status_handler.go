package commands

import (
	"context"

	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"
)

// AssignQuestCommandHandler defines the interface for handling AssignQuestCommand.
type AssignQuestCommandHandler interface {
	Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error)
}

var _ AssignQuestCommandHandler = &assignQuestHandler{}

// assignQuestHandler implements AssignQuestCommandHandler.
type assignQuestHandler struct {
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
}

// NewAssignQuestCommandHandler creates a new instance of AssignQuestCommandHandler.
func NewAssignQuestCommandHandler(unitOfWork ports.UnitOfWork, eventPublisher ports.EventPublisher) AssignQuestCommandHandler {
	return &assignQuestHandler{
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
	}
}

// Handle assigns a quest to a user using domain business rules.
func (h *assignQuestHandler) Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error) {
	// Begin transaction
	if err := h.unitOfWork.Begin(ctx); err != nil {
		return AssignQuestResult{}, errs.WrapInfrastructureError("failed to begin quest assignment transaction", err)
	}

	// Get quest - if not found → 404
	q, err := h.unitOfWork.QuestRepository().GetByID(ctx, cmd.ID)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return AssignQuestResult{}, errs.NewNotFoundErrorWithCause("quest", cmd.ID.String(), err)
	}

	// Use domain logic - business rules errors → 400
	if err := q.AssignTo(cmd.UserID); err != nil {
		_ = h.unitOfWork.Rollback()
		return AssignQuestResult{}, errs.NewDomainValidationErrorWithCause("assignment", "failed to assign quest", err)
	}

	// Save quest - infrastructure error → 500
	if err := h.unitOfWork.QuestRepository().Save(ctx, q); err != nil {
		_ = h.unitOfWork.Rollback()
		return AssignQuestResult{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Publish domain events within the same transaction
	if h.eventPublisher != nil {
		if err := h.eventPublisher.Publish(ctx, q.GetDomainEvents()...); err != nil {
			_ = h.unitOfWork.Rollback()
			return AssignQuestResult{}, errs.WrapInfrastructureError("failed to publish events", err)
		}
	}

	// Commit transaction
	err = h.unitOfWork.Commit(ctx)
	if err != nil {
		return AssignQuestResult{}, errs.WrapInfrastructureError("failed to commit quest assignment transaction", err)
	}

	// Clear events after successful commit
	q.ClearDomainEvents()

	return AssignQuestResult{
		ID:       q.ID(),
		Assignee: cmd.UserID,
		Status:   string(q.Status),
	}, nil
}
