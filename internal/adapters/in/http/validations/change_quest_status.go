package validations

import (
	"quest-manager/internal/generated/servers"

	"github.com/google/uuid"
)

// ValidatedChangeQuestStatusData contains validated data for quest status change
type ValidatedChangeQuestStatusData struct {
	QuestID uuid.UUID
	Status  string
}

// ValidateChangeQuestStatusRequest validates quest status change request
func ValidateChangeQuestStatusRequest(req *servers.ChangeStatusRequest, questIdParam string) (*ValidatedChangeQuestStatusData, *ValidationError) {
	// Validate body
	if err := ValidateBody(req, "body"); err != nil {
		return nil, err
	}

	// Validate Status - only check that it's not empty string
	statusStr := string(req.Status)
	if err := ValidateNotEmpty(statusStr, "status"); err != nil {
		return nil, err
	}

	// Validate QuestId format (UUID)
	questID, err := ValidateUUID(questIdParam, "questId")
	if err != nil {
		return nil, err
	}

	return &ValidatedChangeQuestStatusData{
		QuestID: questID,
		Status:  statusStr,
	}, nil
}
