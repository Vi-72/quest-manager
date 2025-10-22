package commands

import (
	"context"

	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/errs"

	"github.com/google/uuid"
)

// CreateQuestCommandHandler defines a handler for creating a new quest.
type CreateQuestCommandHandler interface {
	Handle(ctx context.Context, cmd CreateQuestCommand) (quest.Quest, error)
}

var _ CreateQuestCommandHandler = &createQuestHandler{}

type createQuestHandler struct {
	executor *CommandExecutor
}

// NewCreateQuestCommandHandler creates a new instance of CreateQuestCommandHandler.
func NewCreateQuestCommandHandler(uowFactory ports.UnitOfWorkFactory, eventPublisher ports.EventPublisher) CreateQuestCommandHandler {
	return &createQuestHandler{
		executor: NewCommandExecutor(uowFactory, eventPublisher),
	}
}

func (h *createQuestHandler) Handle(ctx context.Context, cmd CreateQuestCommand) (quest.Quest, error) {
	var createdQuest quest.Quest

	err := h.executor.Execute(ctx, func(ctx context.Context, uow ports.UnitOfWork) ([]ddd.AggregateRoot, error) {
		var targetLocationID *uuid.UUID
		var executionLocationID *uuid.UUID

		// Create or find target location
		targetLoc, err := location.NewLocation(
			cmd.TargetLocation,
			cmd.TargetAddress,
		)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to create target location", err)
		}

		// Save target location
		err = uow.LocationRepository().Save(ctx, targetLoc)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to save target location", err)
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
				return nil, errs.WrapInfrastructureError("failed to create execution location", err)
			}

			// Save execution location
			err = uow.LocationRepository().Save(ctx, executionLoc)
			if err != nil {
				return nil, errs.WrapInfrastructureError("failed to save execution location", err)
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
			return nil, errs.NewDomainValidationErrorWithCause("quest", "invalid quest data", err)
		}

		// Link quest with created locations
		q.TargetLocationID = targetLocationID
		q.ExecutionLocationID = executionLocationID

		// Save quest
		err = uow.QuestRepository().Save(ctx, q)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to save quest", err)
		}

		createdQuest = q

		// Return aggregates for event publishing
		if executionLoc != targetLoc {
			return []ddd.AggregateRoot{&q, targetLoc, executionLoc}, nil
		}
		return []ddd.AggregateRoot{&q, targetLoc}, nil
	})

	return createdQuest, err
}
