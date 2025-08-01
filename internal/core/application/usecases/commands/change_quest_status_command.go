package commands

import (
	"quest-manager/internal/core/domain/model/quest"

	"github.com/google/uuid"
)

type ChangeQuestStatusCommand struct {
	QuestID uuid.UUID
	Status  quest.Status
}
