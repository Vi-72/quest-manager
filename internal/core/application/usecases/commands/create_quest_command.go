package commands

import (
	"time"

	"github.com/google/uuid"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

type CreateQuestCommand struct {
	Title             string
	Description       string
	Difficulty        quest.Difficulty
	Reward            string
	TargetLocation    kernel.GeoCoordinate
	ExecutionLocation kernel.GeoCoordinate
	Equipment         []string
	Skills            []string
	Creator           string
}

type CreateQuestResult struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Status    quest.Status
}
