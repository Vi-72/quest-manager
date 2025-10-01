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
	unitOfWorkFactory ports.UnitOfWorkFactory
}

// NewCreateQuestCommandHandler creates a new instance of CreateQuestCommandHandler.
func NewCreateQuestCommandHandler(factory ports.UnitOfWorkFactory) CreateQuestCommandHandler {
	return &createQuestHandler{
		unitOfWorkFactory: factory,
	}
}

func (h *createQuestHandler) Handle(ctx context.Context, cmd CreateQuestCommand) (quest.Quest, error) {
	unitOfWork, eventPublisher, err := h.unitOfWorkFactory()
	if err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to create quest unit of work", err)
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

	if err := unitOfWork.Begin(ctx); err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to begin quest creation transaction", err)
	}
	txStarted = true

	var targetLocationID *uuid.UUID
	var executionLocationID *uuid.UUID

	// Create or find target location
	targetLoc, err := location.NewLocation(
		cmd.TargetLocation,
		cmd.TargetAddress,
	)
	if err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to create target location", err)
	}

	// Save target location
	err = unitOfWork.LocationRepository().Save(ctx, targetLoc)
	if err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to save target location", err)
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
			return quest.Quest{}, errs.WrapInfrastructureError("failed to create execution location", err)
		}

		// Save execution location
		err = unitOfWork.LocationRepository().Save(ctx, executionLoc)
		if err != nil {
			return quest.Quest{}, errs.WrapInfrastructureError("failed to save execution location", err)
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
		return quest.Quest{}, errs.NewDomainValidationErrorWithCause("quest", "invalid quest data", err)
	}

	// Link quest with created locations
	q.TargetLocationID = targetLocationID
	q.ExecutionLocationID = executionLocationID

	// Save quest
	err = unitOfWork.QuestRepository().Save(ctx, q)
	if err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Commit transaction
	if err := unitOfWork.Commit(ctx); err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to commit quest creation transaction", err)
	}
	committed = true

	// Publish all domain events asynchronously after successful commit
	if executionLoc != targetLoc {
		PublishDomainEventsAsync(context.Background(), eventPublisher, q, targetLoc, executionLoc)
	} else {
		PublishDomainEventsAsync(context.Background(), eventPublisher, q, targetLoc)
	}

	return q, nil
}
