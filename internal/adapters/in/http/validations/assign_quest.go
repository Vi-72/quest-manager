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
	if err := ValidateBody(req, "body"); err != nil {
		return nil, err
	}

	// Валидация UserId
	userID, err := ValidateUUID(req.UserId, "userId")
	if err != nil {
		return nil, err
	}

	// Валидация QuestId format (UUID)
	questID, err := ValidateUUID(questIdParam, "questId")
	if err != nil {
		return nil, err
	}

	return &ValidatedAssignQuestData{
		QuestID: questID,
		UserID:  userID.String(), // Сохраняем как строку для совместимости
	}, nil
}
