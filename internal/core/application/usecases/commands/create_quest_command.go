package commands

import (
	"quest-manager/internal/core/domain/model/kernel"
)

type CreateQuestCommand struct {
	Title             string
	Description       string
	Difficulty        string // Changed to string, validation in domain
	Reward            int
	DurationMinutes   int
	TargetLocation    kernel.GeoCoordinate
	TargetAddress     *string
	ExecutionLocation kernel.GeoCoordinate
	ExecutionAddress  *string
	Equipment         []string
	Skills            []string
	Creator           string
}
