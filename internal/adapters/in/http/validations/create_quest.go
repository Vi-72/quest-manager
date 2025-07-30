package validations

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
	"strings"
)

// ValidatedCreateQuestData содержит валидированные и обработанные данные
type ValidatedCreateQuestData struct {
	Title             string
	Description       string
	Difficulty        string
	Reward            string
	TargetLocation    kernel.GeoCoordinate
	ExecutionLocation kernel.GeoCoordinate
	Equipment         []string
	Skills            []string
}

// ValidateCreateQuestRequest валидирует технические аспекты запроса (формат, диапазоны, не пустые значения)
func ValidateCreateQuestRequest(req *servers.CreateQuestRequest) (*ValidatedCreateQuestData, *ValidationError) {
	// Валидация body
	if err := ValidateBody(req, "body"); err != nil {
		return nil, err
	}

	// Техническая валидация title
	title, err := TrimAndValidateString(req.Title, "title")
	if err != nil {
		return nil, err
	}

	// Техническая валидация description
	description, err := TrimAndValidateString(req.Description, "description")
	if err != nil {
		return nil, err
	}

	// Техническая валидация difficulty (только не пустое)
	if err := ValidateNotEmpty(string(req.Difficulty), "difficulty"); err != nil {
		return nil, err
	}

	// Валидация и конвертация target_location
	targetLocation, err := ConvertAPICoordinateToKernel(req.TargetLocation)
	if err != nil {
		return nil, NewValidationErrorWithCause("target_location", err.Message, err.Cause)
	}

	// Валидация и конвертация execution_location
	executionLocation, err := ConvertAPICoordinateToKernel(req.ExecutionLocation)
	if err != nil {
		return nil, NewValidationErrorWithCause("execution_location", err.Message, err.Cause)
	}

	// Обработка опциональных полей
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
		reward = strings.TrimSpace(*req.Reward)
	}

	return &ValidatedCreateQuestData{
		Title:             title,
		Description:       description,
		Difficulty:        string(req.Difficulty), // Передаем как есть, домен проверит
		Reward:            reward,
		TargetLocation:    targetLocation,
		ExecutionLocation: executionLocation,
		Equipment:         equipment,
		Skills:            skills,
	}, nil
}
