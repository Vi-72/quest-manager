package validations

import (
	v1 "quest-manager/api/http/quests/v1"

	"github.com/google/uuid"
)

// ValidatedChangeQuestStatusData contains validated data for quest status change
type ValidatedChangeQuestStatusData struct {
	QuestID uuid.UUID
	Status  string
}

// ValidateChangeQuestStatusRequest validates and converts quest status change request
// Note: Status enum and UUID format validations are now handled by OpenAPI
func ValidateChangeQuestStatusRequest(req *v1.ChangeStatusRequest, questIdParam string) (*ValidatedChangeQuestStatusData, *ValidationError) {
	// Validate body (still needed for nil check)
	if err := ValidateBody(req, "body"); err != nil {
		return nil, err
	}

	// OpenAPI now handles:
	// - status: enum validation (created, posted, assigned, in_progress, declined, completed)
	// - quest_id: format uuid validation (path parameter)

	// Convert OpenAPI types to Go types
	statusStr := string(req.Status)          // Already validated by OpenAPI enum
	questID, err := uuid.Parse(questIdParam) // questIdParam already validated by OpenAPI
	if err != nil {
		// This should never happen since OpenAPI validates format, but keep as safety net
		return nil, NewValidationErrorWithCause("questId", "invalid UUID format", err)
	}

	return &ValidatedChangeQuestStatusData{
		QuestID: questID,
		Status:  statusStr,
	}, nil
}
