package http

import (
	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/domain/model/quest"
)

// QuestToAPI converts domain quest to API format
func QuestToAPI(q quest.Quest) v1.Quest {
	// Convert target and execution locations
	// Note: addresses are not denormalized in Quest, they're in separate Location entities
	targetLocation := validations.ConvertKernelCoordinateToAPI(q.TargetLocation, nil)
	executionLocation := validations.ConvertKernelCoordinateToAPI(q.ExecutionLocation, nil)

	// Convert optional fields to pointers, ensure empty arrays instead of null
	equipment := &q.Equipment

	skills := &q.Skills

	// Convert location IDs to strings if present
	var targetLocationId *string
	if q.TargetLocationID != nil {
		id := q.TargetLocationID.String()
		targetLocationId = &id
	}

	var executionLocationId *string
	if q.ExecutionLocationID != nil {
		id := q.ExecutionLocationID.String()
		executionLocationId = &id
	}

	return v1.Quest{
		Id:                  q.ID().String(),
		Title:               q.Title,
		Description:         q.Description,
		Difficulty:          v1.QuestDifficulty(q.Difficulty),
		Reward:              q.Reward,
		DurationMinutes:     q.DurationMinutes,
		TargetLocation:      targetLocation,
		ExecutionLocation:   executionLocation,
		Equipment:           equipment,
		Skills:              skills,
		Status:              v1.QuestStatus(q.Status),
		Creator:             q.Creator,
		Assignee:            q.Assignee,
		CreatedAt:           q.CreatedAt,
		UpdatedAt:           q.UpdatedAt,
		TargetLocationId:    targetLocationId,
		ExecutionLocationId: executionLocationId,
	}
}
