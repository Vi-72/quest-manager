package questrepo

import (
	"strings"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"

	"github.com/google/uuid"
)

// DomainToDTO converts Quest domain model to QuestDTO for DB.
func DomainToDTO(q quest.Quest) QuestDTO {
	return QuestDTO{
		ID:                 q.ID.String(),
		Title:              q.Title,
		Description:        q.Description,
		Difficulty:         string(q.Difficulty),
		Reward:             q.Reward,
		TargetLatitude:     q.TargetLocation.Latitude(),
		TargetLongitude:    q.TargetLocation.Longitude(),
		ExecutionLatitude:  q.ExecutionLocation.Latitude(),
		ExecutionLongitude: q.ExecutionLocation.Longitude(),
		Equipment:          strings.Join(q.Equipment, ","),
		Skills:             strings.Join(q.Skills, ","),
		Status:             string(q.Status),
		Creator:            q.Creator,
		Assignee:           q.Assignee,
		CreatedAt:          q.CreatedAt,
		UpdatedAt:          q.UpdatedAt,
	}
}

// DtoToDomain converts QuestDTO to domain model Quest.
func DtoToDomain(dto QuestDTO) (quest.Quest, error) {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return quest.Quest{}, err
	}

	targetCoord, err := kernel.NewGeoCoordinate(dto.TargetLatitude, dto.TargetLongitude)
	if err != nil {
		return quest.Quest{}, err
	}

	execCoord, err := kernel.NewGeoCoordinate(dto.ExecutionLatitude, dto.ExecutionLongitude)
	if err != nil {
		return quest.Quest{}, err
	}

	var equipment []string
	if dto.Equipment != "" {
		equipment = strings.Split(dto.Equipment, ",")
	}

	var skills []string
	if dto.Skills != "" {
		skills = strings.Split(dto.Skills, ",")
	}

	return quest.Quest{
		ID:                id,
		Title:             dto.Title,
		Description:       dto.Description,
		Difficulty:        quest.Difficulty(dto.Difficulty),
		Reward:            dto.Reward,
		TargetLocation:    targetCoord,
		ExecutionLocation: execCoord,
		Equipment:         equipment,
		Skills:            skills,
		Status:            quest.Status(dto.Status),
		Creator:           dto.Creator,
		Assignee:          dto.Assignee,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}, nil
}
