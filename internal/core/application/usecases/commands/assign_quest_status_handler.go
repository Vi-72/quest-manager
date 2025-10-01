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
	unitOfWorkFactory ports.UnitOfWorkFactory
}

// NewAssignQuestCommandHandler creates a new instance of AssignQuestCommandHandler.
func NewAssignQuestCommandHandler(factory ports.UnitOfWorkFactory) AssignQuestCommandHandler {
	return &assignQuestHandler{
		unitOfWorkFactory: factory,
	}
}

// Handle assigns a quest to a user using domain business rules.
func (h *assignQuestHandler) Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error) {
	unitOfWork, eventPublisher, err := h.unitOfWorkFactory()
	if err != nil {
		return AssignQuestResult{}, errs.WrapInfrastructureError("failed to create quest assignment unit of work", err)
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
		return AssignQuestResult{}, errs.WrapInfrastructureError("failed to begin quest assignment transaction", err)
	}
	txStarted = true

	// Get quest - if not found → 404
	q, err := unitOfWork.QuestRepository().GetByID(ctx, cmd.ID)
	if err != nil {
		return AssignQuestResult{}, errs.NewNotFoundErrorWithCause("quest", cmd.ID.String(), err)
	}

	// Use domain logic - business rules errors → 400
	if err := q.AssignTo(cmd.UserID); err != nil {
		return AssignQuestResult{}, errs.NewDomainValidationErrorWithCause("assignment", "failed to assign quest", err)
	}

	// Save quest - infrastructure error → 500
	if err := unitOfWork.QuestRepository().Save(ctx, q); err != nil {
		return AssignQuestResult{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Publish domain events within the same transaction
	if eventPublisher != nil {
		if err := eventPublisher.Publish(ctx, q.GetDomainEvents()...); err != nil {
			return AssignQuestResult{}, errs.WrapInfrastructureError("failed to publish events", err)
		}
	}

	// Commit transaction
	if err := unitOfWork.Commit(ctx); err != nil {
		return AssignQuestResult{}, errs.WrapInfrastructureError("failed to commit quest assignment transaction", err)
	}
	committed = true

	// Clear events after successful commit
	q.ClearDomainEvents()

	return AssignQuestResult{
		ID:       q.ID(),
		Assignee: cmd.UserID,
		Status:   string(q.Status),
	}, nil
}
