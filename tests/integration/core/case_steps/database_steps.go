package casesteps

import (
	"context"

	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
)

// CountInitialDatabaseRecords is a helper to get initial record counts before test operations
func CountInitialDatabaseRecords(
	ctx context.Context,
	questRepo interface {
		FindAll(ctx context.Context) ([]quest.Quest, error)
	},
	locationRepo interface {
		FindAll(ctx context.Context) ([]*location.Location, error)
	},
) (initialQuests []quest.Quest, initialLocations []*location.Location, err error) {
	initialQuests, err = questRepo.FindAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	initialLocations, err = locationRepo.FindAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	return initialQuests, initialLocations, nil
}
