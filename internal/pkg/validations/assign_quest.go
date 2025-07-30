package validations

import (
	"quest-manager/internal/generated/servers"

	"github.com/google/uuid"
)

// ValidatedAssignQuestData содержит валидированные данные для назначения квеста
type ValidatedAssignQuestData struct {
	QuestID uuid.UUID
	UserID  string
}

// ValidateAssignQuestRequest валидирует запрос назначения квеста
func ValidateAssignQuestRequest(req *servers.AssignQuestRequest, questIdParam string) (*ValidatedAssignQuestData, *ValidationError) {
	// Валидация body
	if req == nil {
		return nil, NewValidationError("body", "is required")
	}

	// Валидация UserId
	if req.UserId == "" {
		return nil, NewValidationError("userId", "is required")
	}

	// Валидация UserId format (UUID)
	_, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, NewValidationErrorWithCause("userId", "must be a valid UUID format (e.g. d5cde057-d462-419b-9428-42eebe22a85e)", err)
	}

	// Валидация QuestId format (UUID)
	questID, err := uuid.Parse(questIdParam)
	if err != nil {
		return nil, NewValidationErrorWithCause("questId", "must be a valid UUID format", err)
	}

	return &ValidatedAssignQuestData{
		QuestID: questID,
		UserID:  req.UserId,
	}, nil
}
