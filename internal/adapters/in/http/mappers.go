package http

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
)

// Coordinate mappers
func APICoordinateToKernel(apiCoord servers.Coordinate) (kernel.GeoCoordinate, error) {
	return kernel.NewGeoCoordinate(float64(apiCoord.Latitude), float64(apiCoord.Longitude))
}

func KernelCoordinateToAPI(kernelCoord kernel.GeoCoordinate) servers.Coordinate {
	return servers.Coordinate{
		Latitude:  float32(kernelCoord.Latitude()),
		Longitude: float32(kernelCoord.Longitude()),
	}
}

// Quest mappers
func QuestToAPI(q quest.Quest) servers.Quest {
	return servers.Quest{
		Id:                string(q.ID.String()),
		Title:             q.Title,
		Description:       q.Description,
		Difficulty:        servers.QuestDifficulty(q.Difficulty),
		Reward:            &q.Reward,
		TargetLocation:    KernelCoordinateToAPI(q.TargetLocation),
		ExecutionLocation: KernelCoordinateToAPI(q.ExecutionLocation),
		Equipment:         &q.Equipment,
		Skills:            &q.Skills,
		Status:            servers.QuestStatus(q.Status),
		Creator:           q.Creator,
		Assignee:          q.Assignee,
		CreatedAt:         q.CreatedAt,
		UpdatedAt:         q.UpdatedAt,
	}
}

func CreateQuestRequestToDomain(req servers.CreateQuestRequest, creator string) (quest.Quest, error) {
	targetLocation, err := APICoordinateToKernel(req.TargetLocation)
	if err != nil {
		return quest.Quest{}, err
	}

	executionLocation, err := APICoordinateToKernel(req.ExecutionLocation)
	if err != nil {
		return quest.Quest{}, err
	}

	equipment := []string{}
	if req.Equipment != nil {
		equipment = *req.Equipment
	}

	skills := []string{}
	if req.Skills != nil {
		skills = *req.Skills
	}

	reward := ""
	if req.Reward != nil {
		reward = *req.Reward
	}

	return quest.NewQuest(
		req.Title,
		req.Description,
		quest.Difficulty(req.Difficulty),
		reward,
		targetLocation,
		executionLocation,
		creator,
		equipment,
		skills,
	), nil
}
