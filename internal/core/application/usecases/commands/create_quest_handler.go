package commands

import (
	"context"

	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"

	"github.com/google/uuid"
)

// CreateQuestCommandHandler defines a handler for creating a new quest.
type CreateQuestCommandHandler interface {
	Handle(ctx context.Context, cmd CreateQuestCommand) (quest.Quest, error)
}

var _ CreateQuestCommandHandler = &createQuestHandler{}

type createQuestHandler struct {
	txManager ports.TransactionManager
}

// NewCreateQuestCommandHandler creates a new instance of CreateQuestCommandHandler.
func NewCreateQuestCommandHandler(txManager ports.TransactionManager) CreateQuestCommandHandler {
	return &createQuestHandler{
		txManager: txManager,
	}
}

func (h *createQuestHandler) Handle(ctx context.Context, cmd CreateQuestCommand) (quest.Quest, error) {
	var createdQuest quest.Quest

	err := h.txManager.RunInTransaction(ctx, func(ctx context.Context, repos ports.Repositories) error {
		var targetLocationID *uuid.UUID
		var executionLocationID *uuid.UUID

		// Create or find target location
		targetLoc, err := location.NewLocation(
			cmd.TargetLocation,
			cmd.TargetAddress,
		)
		if err != nil {
			return errs.WrapInfrastructureError("failed to create target location", err)
		}

		// Save target location
		err = repos.Location.Save(ctx, targetLoc)
		if err != nil {
			return errs.WrapInfrastructureError("failed to save target location", err)
		}
		targetLocID := targetLoc.ID()
		targetLocationID = &targetLocID

		// Create or find execution location (can be the same as target)
		var executionLoc *location.Location
		if cmd.TargetLocation.Equals(cmd.ExecutionLocation) {
			executionLoc = targetLoc
			executionLocationID = targetLocationID
		} else {
			executionLoc, err = location.NewLocation(
				cmd.ExecutionLocation,
				cmd.ExecutionAddress,
			)
			if err != nil {
				return errs.WrapInfrastructureError("failed to create execution location", err)
			}

			// Save execution location
			err = repos.Location.Save(ctx, executionLoc)
			if err != nil {
				return errs.WrapInfrastructureError("failed to save execution location", err)
			}
			executionLocID := executionLoc.ID()
			executionLocationID = &executionLocID
		}

		// Create quest
		q, err := quest.NewQuest(
			cmd.Title,
			cmd.Description,
			cmd.Difficulty,
			cmd.Reward,
			cmd.DurationMinutes,
			cmd.TargetLocation,
			cmd.ExecutionLocation,
			cmd.Creator,
			cmd.Equipment,
			cmd.Skills,
		)
		if err != nil {
			return errs.NewDomainValidationErrorWithCause("quest", "invalid quest data", err)
		}

		// Link quest with created locations
		q.TargetLocationID = targetLocationID
		q.ExecutionLocationID = executionLocationID

		// Save quest
		err = repos.Quest.Save(ctx, q)
		if err != nil {
			return errs.WrapInfrastructureError("failed to save quest", err)
		}

		// Publish events synchronously in same transaction
		events := q.GetDomainEvents()
		events = append(events, targetLoc.GetDomainEvents()...)
		// Only add execution location events if it's different from target location
		if executionLoc != targetLoc {
			events = append(events, executionLoc.GetDomainEvents()...)
		}
		if err := repos.Event.Publish(ctx, events...); err != nil {
			return errs.WrapInfrastructureError("failed to publish events", err)
		}

		// Clear events after successful publication
		q.ClearDomainEvents()
		targetLoc.ClearDomainEvents()
		if executionLoc != nil && executionLoc != targetLoc {
			executionLoc.ClearDomainEvents()
		}

		createdQuest = q
		return nil
	})

	return createdQuest, err
}
