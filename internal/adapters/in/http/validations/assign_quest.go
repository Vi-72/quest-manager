package validations

import (
	"quest-manager/internal/generated/servers"

	"github.com/google/uuid"
)

// ValidatedAssignQuestData contains validated data for quest assignment
type ValidatedAssignQuestData struct {
	QuestID uuid.UUID
	UserID  string
}

// ValidateAssignQuestRequest validates quest assignment request
func ValidateAssignQuestRequest(req *servers.AssignQuestRequest, questIdParam string) (*ValidatedAssignQuestData, *ValidationError) {
	// Validate body
	if err := ValidateBody(req, "body"); err != nil {
		return nil, err
	}

	// Validate UserId
	userID, err := ValidateUUID(req.UserId, "userId")
	if err != nil {
		return nil, err
	}

	// Validate QuestId format (UUID)
	questID, err := ValidateUUID(questIdParam, "questId")
	if err != nil {
		return nil, err
	}

	return &ValidatedAssignQuestData{
		QuestID: questID,
		UserID:  userID.String(), // Store as string for compatibility
	}, nil
}
