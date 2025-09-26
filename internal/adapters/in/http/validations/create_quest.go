package validations

import (
	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/core/domain/model/kernel"
)

// ValidatedCreateQuestData contains validated and processed data
type ValidatedCreateQuestData struct {
	Title             string
	Description       string
	Difficulty        string
	Reward            int
	DurationMinutes   int
	TargetLocation    kernel.GeoCoordinate
	TargetAddress     *string
	ExecutionLocation kernel.GeoCoordinate
	ExecutionAddress  *string
	Equipment         []string
	Skills            []string
}

// ValidateCreateQuestRequest validates and converts request data
// Note: Basic validations (required fields, formats, ranges) are now handled by OpenAPI
func ValidateCreateQuestRequest(req *v1.CreateQuestRequest) (*ValidatedCreateQuestData, *ValidationError) {
	// Validate body (still needed for nil check)
	if err := ValidateBody(req, "body"); err != nil {
		return nil, err
	}

	// OpenAPI now handles:
	// - title: minLength, maxLength, pattern (no whitespace-only)
	// - description: minLength, maxLength, pattern (no whitespace-only)
	// - difficulty: enum validation
	// - reward: minimum 1, maximum 5
	// - duration_minutes: minimum 1, maximum 10080

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
		Title:             req.Title,              // OpenAPI validates minLength, maxLength, pattern
		Description:       req.Description,        // OpenAPI validates minLength, maxLength, pattern
		Difficulty:        string(req.Difficulty), // OpenAPI validates enum
		Reward:            req.Reward,             // OpenAPI validates minimum 1, maximum 5
		DurationMinutes:   req.DurationMinutes,    // OpenAPI validates minimum 1, maximum 10080
		TargetLocation:    targetLocation,
		TargetAddress:     req.TargetLocation.Address,
		ExecutionLocation: executionLocation,
		ExecutionAddress:  req.ExecutionLocation.Address,
		Equipment:         equipment,
		Skills:            skills,
	}, nil
}
