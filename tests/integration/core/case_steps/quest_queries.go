package casesteps

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

// GetQuestByIDStep получает квест по ID
func GetQuestByIDStep(
	ctx context.Context,
	handler queries.GetQuestByIDQueryHandler,
	questID uuid.UUID,
) (quest.Quest, error) {
	return handler.Handle(ctx, questID)
}

// ListQuestsStep получает список квестов
func ListQuestsStep(
	ctx context.Context,
	handler queries.ListQuestsQueryHandler,
	status *quest.Status,
) ([]quest.Quest, error) {
	return handler.Handle(ctx, status)
}

// ListAssignedQuestsStep получает список назначенных квестов для пользователя
func ListAssignedQuestsStep(
	ctx context.Context,
	handler queries.ListAssignedQuestsQueryHandler,
	userID string,
) ([]quest.Quest, error) {
	return handler.Handle(ctx, userID)
}

// SearchQuestsByRadiusStep ищет квесты в радиусе
func SearchQuestsByRadiusStep(
	ctx context.Context,
	handler queries.SearchQuestsByRadiusQueryHandler,
	lat, lon, radiusKm float64,
) ([]quest.Quest, error) {
	center := kernel.GeoCoordinate{
		Lat: lat,
		Lon: lon,
	}

	return handler.Handle(ctx, center, radiusKm)
}
