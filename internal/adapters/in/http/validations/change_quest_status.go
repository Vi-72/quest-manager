package validations

import (
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"

	"github.com/google/uuid"
)

// validQuestStatuses содержит все допустимые значения статуса квеста для изменения
var validQuestStatuses = []string{
	string(quest.StatusCreated),
	string(quest.StatusPosted),
	string(quest.StatusAssigned),
	string(quest.StatusInProgress),
	string(quest.StatusDeclined),
	string(quest.StatusCompleted),
}

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

	// Валидация Status
	statusStr := string(req.Status)
	if err := ValidateNotEmpty(statusStr, "status"); err != nil {
		return nil, err
	}

	// Валидация что статус является допустимым значением
	if err := ValidateEnum(statusStr, validQuestStatuses, "status"); err != nil {
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
