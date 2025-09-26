package validations

import (
	v1 "quest-manager/api/http/quests/v1"

	"github.com/google/uuid"
)

// ValidatedAssignQuestData contains validated data for quest assignment
type ValidatedAssignQuestData struct {
	QuestID uuid.UUID
	UserID  string
}

// ValidateAssignQuestRequest validates and converts quest assignment request
// Note: UUID format validations are now handled by OpenAPI
func ValidateAssignQuestRequest(req *v1.AssignQuestRequest, questIdParam string) (*ValidatedAssignQuestData, *ValidationError) {
	// Validate body (still needed for nil check)
	if err := ValidateBody(req, "body"); err != nil {
		return nil, err
	}

	// OpenAPI now handles:
	// - user_id: format uuid validation
	// - quest_id: format uuid validation (path parameter)

	// Convert OpenAPI UUID types to Go types
	userID := req.UserId
	questID, err := uuid.Parse(questIdParam) // questIdParam already validated by OpenAPI
	if err != nil {
		// This should never happen since OpenAPI validates format, but keep as safety net
		return nil, NewValidationErrorWithCause("questId", "invalid UUID format", err)
	}

	return &ValidatedAssignQuestData{
		QuestID: questID,
		UserID:  userID.String(), // Convert UUID to string for compatibility
	}, nil
}
