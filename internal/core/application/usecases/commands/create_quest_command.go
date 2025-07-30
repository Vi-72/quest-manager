package commands

import (
	"time"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"

	"github.com/google/uuid"
)

type CreateQuestCommand struct {
	Title             string
	Description       string
	Difficulty        string // Изменено на string, валидация в домене
	Reward            string
	TargetLocation    kernel.GeoCoordinate
	ExecutionLocation kernel.GeoCoordinate
	Equipment         []string
	Skills            []string
	Creator           string
}

type CreateQuestResult struct {
	ID                  uuid.UUID
	CreatedAt           time.Time
	Status              quest.Status
	TargetLocationID    *uuid.UUID // ID created location (if any)
	ExecutionLocationID *uuid.UUID // ID created location (if any)
}
