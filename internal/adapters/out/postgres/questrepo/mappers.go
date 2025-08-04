package questrepo

import (
	"strings"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/pkg/ddd"

	"github.com/google/uuid"
)

// DomainToDTO converts Quest domain model to QuestDTO for DB.
func DomainToDTO(q quest.Quest) QuestDTO {
	dto := QuestDTO{
		ID:                 q.ID().String(),
		Title:              q.Title,
		Description:        q.Description,
		Difficulty:         string(q.Difficulty),
		Reward:             q.Reward,
		DurationMinutes:    q.DurationMinutes,
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

	// Опциональные ссылки на локации
	if q.TargetLocationID != nil {
		targetLocationIDStr := q.TargetLocationID.String()
		dto.TargetLocationID = &targetLocationIDStr
	}
	if q.ExecutionLocationID != nil {
		executionLocationIDStr := q.ExecutionLocationID.String()
		dto.ExecutionLocationID = &executionLocationIDStr
	}

	return dto
}

// DtoToDomainWithAddress converts QuestWithAddressDTO to Quest domain model
func DtoToDomainWithAddress(dto QuestWithAddressDTO) (quest.Quest, error) {
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

	q, err := dtoToDomainCommon(dto.QuestDTO, id, targetCoord, execCoord)
	if err != nil {
		return quest.Quest{}, err
	}

	// Add addresses from DTO
	q.TargetAddress = dto.TargetAddress
	q.ExecutionAddress = dto.ExecutionAddress

	return q, nil
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

	return dtoToDomainCommon(dto, id, targetCoord, execCoord)
}

// dtoToDomainCommon contains shared logic for converting DTO to domain
func dtoToDomainCommon(dto QuestDTO, id uuid.UUID, targetCoord, execCoord kernel.GeoCoordinate) (quest.Quest, error) {

	var equipment []string
	if dto.Equipment != "" {
		equipment = strings.Split(dto.Equipment, ",")
	} else {
		equipment = []string{} // Нормализация: всегда возвращаем [], а не nil
	}

	var skills []string
	if dto.Skills != "" {
		skills = strings.Split(dto.Skills, ",")
	} else {
		skills = []string{} // Нормализация: всегда возвращаем [], а не nil
	}

	q := quest.Quest{
		BaseAggregate:     ddd.NewBaseAggregate(id),
		Title:             dto.Title,
		Description:       dto.Description,
		Difficulty:        quest.Difficulty(dto.Difficulty),
		Reward:            dto.Reward,
		DurationMinutes:   dto.DurationMinutes,
		TargetLocation:    targetCoord,
		ExecutionLocation: execCoord,
		Equipment:         equipment,
		Skills:            skills,
		Status:            quest.Status(dto.Status),
		Creator:           dto.Creator,
		Assignee:          dto.Assignee,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}

	// Опциональные ссылки на локации
	if dto.TargetLocationID != nil {
		targetLocationID, err := uuid.Parse(*dto.TargetLocationID)
		if err != nil {
			return quest.Quest{}, err
		}
		q.TargetLocationID = &targetLocationID
	}
	if dto.ExecutionLocationID != nil {
		executionLocationID, err := uuid.Parse(*dto.ExecutionLocationID)
		if err != nil {
			return quest.Quest{}, err
		}
		q.ExecutionLocationID = &executionLocationID
	}

	return q, nil
}
