package commands

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/errs"
)

// ChangeQuestStatusCommandHandler defines the interface for handling quest status changes.
type ChangeQuestStatusCommandHandler interface {
	Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (ChangeQuestStatusResult, error)
}

type changeQuestStatusHandler struct {
	executor *CommandExecutor
}

// NewChangeQuestStatusCommandHandler creates a new ChangeQuestStatusCommandHandler instance.
func NewChangeQuestStatusCommandHandler(uowFactory ports.UnitOfWorkFactory, eventPublisher ports.EventPublisher) ChangeQuestStatusCommandHandler {
	return &changeQuestStatusHandler{
		executor: NewCommandExecutor(uowFactory, eventPublisher),
	}
}

// Handle updates the quest status with validation and domain business rules.
func (h *changeQuestStatusHandler) Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (ChangeQuestStatusResult, error) {
	// Validate status - this is domain validation error → 400
	if !quest.IsValidStatus(string(cmd.Status)) {
		return ChangeQuestStatusResult{}, errs.NewDomainValidationError("status", "must be one of 'created', 'posted', 'assigned', 'in_progress', 'declined', 'completed'")
	}

	var result ChangeQuestStatusResult

	err := h.executor.Execute(ctx, func(ctx context.Context, uow ports.UnitOfWork) ([]ddd.AggregateRoot, error) {
		// Get quest - if not found → 404
		q, err := uow.QuestRepository().GetByID(ctx, cmd.QuestID)
		if err != nil {
			return nil, errs.NewNotFoundErrorWithCause("quest", cmd.QuestID.String(), err)
		}

		// Use domain logic for status change - domain validation error → 400
		if err := q.ChangeStatus(cmd.Status); err != nil {
			return nil, errs.NewDomainValidationErrorWithCause("status", "invalid status transition", err)
		}

		// Save quest - infrastructure error → 500
		if err := uow.QuestRepository().Save(ctx, q); err != nil {
			return nil, errs.WrapInfrastructureError("failed to save quest", err)
		}

		result = ChangeQuestStatusResult{
			ID:       q.ID(),
			Assignee: q.Assignee, // Now both are *uuid.UUID
			Status:   string(q.Status),
		}

		// Return aggregate for event publishing
		return []ddd.AggregateRoot{&q}, nil
	})

	return result, err
}
