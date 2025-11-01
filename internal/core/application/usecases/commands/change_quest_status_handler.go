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
	txManager ports.TransactionManager
}

// NewChangeQuestStatusCommandHandler creates a new ChangeQuestStatusCommandHandler instance.
func NewChangeQuestStatusCommandHandler(txManager ports.TransactionManager) ChangeQuestStatusCommandHandler {
	return &changeQuestStatusHandler{
		txManager: txManager,
	}
}

// Handle updates the quest status with validation and domain business rules.
func (h *changeQuestStatusHandler) Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (ChangeQuestStatusResult, error) {
	// Validate status - this is domain validation error → 400
	if !quest.IsValidStatus(string(cmd.Status)) {
		return ChangeQuestStatusResult{}, errs.NewDomainValidationError("status", "must be one of 'created', 'posted', 'assigned', 'in_progress', 'declined', 'completed'")
	}

	var result ChangeQuestStatusResult

	err := h.txManager.RunInTransaction(ctx, func(ctx context.Context, repos ports.Repositories) error {
		// Get quest - if not found → 404
		q, err := repos.Quest.GetByID(ctx, cmd.QuestID)
		if err != nil {
			return errs.NewNotFoundErrorWithCause("quest", cmd.QuestID.String(), err)
		}

		// Use domain logic for status change - domain validation error → 400
		if err := q.ChangeStatus(cmd.Status); err != nil {
			return errs.NewDomainValidationErrorWithCause("status", "invalid status transition", err)
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

		result = ChangeQuestStatusResult{
			ID:       q.ID(),
			Assignee: q.Assignee, // Now both are *uuid.UUID
			Status:   string(q.Status),
		}

		return nil
	})

	return result, err
}
