package commands

import (
	"quest-manager/internal/core/domain/model/kernel"
)

type CreateQuestCommand struct {
	Title             string
	Description       string
	Difficulty        string // Изменено на string, валидация в домене
	Reward            int
	DurationMinutes   int
	TargetLocation    kernel.GeoCoordinate
	ExecutionLocation kernel.GeoCoordinate
	Equipment         []string
	Skills            []string
	Creator           string
}
