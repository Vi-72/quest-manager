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
	txManager ports.TransactionManager
}

// NewAssignQuestCommandHandler creates a new instance of AssignQuestCommandHandler.
func NewAssignQuestCommandHandler(txManager ports.TransactionManager) AssignQuestCommandHandler {
	return &assignQuestHandler{
		txManager: txManager,
	}
}

// Handle assigns a quest to a user using domain business rules.
func (h *assignQuestHandler) Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error) {
	var result AssignQuestResult

	err := h.txManager.RunInTransaction(ctx, func(ctx context.Context, repos ports.Repositories) error {
		// Get quest - if not found → 404
		q, err := repos.Quest.GetByID(ctx, cmd.ID)
		if err != nil {
			return errs.NewNotFoundErrorWithCause("quest", cmd.ID.String(), err)
		}

		// Use domain logic - business rules errors → 400
		if err := q.AssignTo(cmd.UserID); err != nil {
			return errs.NewDomainValidationErrorWithCause("assignment", "failed to assign quest", err)
		}

		// Save quest - infrastructure error → 500
		if err := repos.Quest.Save(ctx, q); err != nil {
			return errs.WrapInfrastructureError("failed to save quest", err)
		}

		// Publish events synchronously in same transaction
		if err := repos.Event.Publish(ctx, q.GetDomainEvents()...); err != nil {
			return errs.WrapInfrastructureError("failed to publish events", err)
		}

		// Clear events after successful publication
		q.ClearDomainEvents()

		result = AssignQuestResult{
			ID:       q.ID(),
			Assignee: cmd.UserID,
			Status:   string(q.Status),
		}

		return nil
	})

	return result, err
}
