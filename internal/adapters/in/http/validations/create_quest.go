package validations

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
)

// ValidatedCreateQuestData contains validated and processed data
type ValidatedCreateQuestData struct {
	Title             string
	Description       string
	Difficulty        string
	Reward            int
	DurationMinutes   int
	TargetLocation    kernel.GeoCoordinate
	ExecutionLocation kernel.GeoCoordinate
	Equipment         []string
	Skills            []string
}

// ValidateCreateQuestRequest validates technical aspects of request (format, ranges, non-empty values)
func ValidateCreateQuestRequest(req *servers.CreateQuestRequest) (*ValidatedCreateQuestData, *ValidationError) {
	// Validate body
	if err := ValidateBody(req, "body"); err != nil {
		return nil, err
	}

	// Technical validation of title
	title, err := TrimAndValidateString(req.Title, "title")
	if err != nil {
		return nil, err
	}

	// Technical validation of description
	description, err := TrimAndValidateString(req.Description, "description")
	if err != nil {
		return nil, err
	}

	// Technical validation of difficulty (only non-empty)
	if err := ValidateNotEmpty(string(req.Difficulty), "difficulty"); err != nil {
		return nil, err
	}

	// Basic technical validation of duration_minutes (only positive number)
	if req.DurationMinutes <= 0 {
		return nil, NewValidationError("duration_minutes", "must be a positive number")
	}

	// Basic technical validation of reward (only positive number)
	if req.Reward <= 0 {
		return nil, NewValidationError("reward", "must be a positive number")
	}

	// Validate and convert target_location
	targetLocation, err := ConvertAPICoordinateToKernel(req.TargetLocation)
	if err != nil {
		return nil, NewValidationErrorWithCause("target_location", err.Message, err.Cause)
	}

	// Validate and convert execution_location
	executionLocation, err := ConvertAPICoordinateToKernel(req.ExecutionLocation)
	if err != nil {
		return nil, NewValidationErrorWithCause("execution_location", err.Message, err.Cause)
	}

	// Process optional fields
	equipment := []string{}
	if req.Equipment != nil {
		equipment = *req.Equipment
	}

	skills := []string{}
	if req.Skills != nil {
		skills = *req.Skills
	}

	return &ValidatedCreateQuestData{
		Title:             title,
		Description:       description,
		Difficulty:        string(req.Difficulty), // Pass as is, domain will validate
		Reward:            req.Reward,
		DurationMinutes:   req.DurationMinutes,
		TargetLocation:    targetLocation,
		ExecutionLocation: executionLocation,
		Equipment:         equipment,
		Skills:            skills,
	}, nil
}
