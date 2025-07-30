package validations

import (
	"quest-manager/internal/generated/servers"

	"github.com/google/uuid"
)

// ValidatedChangeQuestStatusData содержит валидированные данные для изменения статуса квеста
type ValidatedChangeQuestStatusData struct {
	QuestID uuid.UUID
	Status  string
}

// ValidateChangeQuestStatusRequest валидирует запрос изменения статуса квеста
func ValidateChangeQuestStatusRequest(req *servers.ChangeStatusRequest, questIdParam string) (*ValidatedChangeQuestStatusData, *ValidationError) {
	// Валидация body
	if err := ValidateBody(req, "body"); err != nil {
		return nil, err
	}

	// Валидация Status - только проверяем что не пустая строка
	statusStr := string(req.Status)
	if err := ValidateNotEmpty(statusStr, "status"); err != nil {
		return nil, err
	}

	// Валидация QuestId format (UUID)
	questID, err := ValidateUUID(questIdParam, "questId")
	if err != nil {
		return nil, err
	}

	return &ValidatedChangeQuestStatusData{
		QuestID: questID,
		Status:  statusStr,
	}, nil
}
