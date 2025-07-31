package http

import (
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
)

// QuestToAPI конвертирует доменный квест в API формат
func QuestToAPI(q quest.Quest) servers.Quest {
	// Convert target and execution locations
	targetLocation := validations.ConvertKernelCoordinateToAPI(q.TargetLocation)
	executionLocation := validations.ConvertKernelCoordinateToAPI(q.ExecutionLocation)

	// Convert optional fields to pointers
	var equipment *[]string
	if len(q.Equipment) > 0 {
		equipment = &q.Equipment
	}

	var skills *[]string
	if len(q.Skills) > 0 {
		skills = &q.Skills
	}

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

	return servers.Quest{
		Id:                  q.ID().String(),
		Title:               q.Title,
		Description:         q.Description,
		Difficulty:          servers.QuestDifficulty(q.Difficulty),
		Reward:              q.Reward,
		DurationMinutes:     q.DurationMinutes,
		TargetLocation:      targetLocation,
		ExecutionLocation:   executionLocation,
		Equipment:           equipment,
		Skills:              skills,
		Status:              servers.QuestStatus(q.Status),
		Creator:             q.Creator,
		Assignee:            q.Assignee,
		CreatedAt:           q.CreatedAt,
		UpdatedAt:           q.UpdatedAt,
		TargetLocationId:    targetLocationId,
		ExecutionLocationId: executionLocationId,
	}
}
