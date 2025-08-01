package casesteps

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// CreateQuestStep выполняет создание квеста
func CreateQuestStep(
	ctx context.Context,
	handler commands.CreateQuestCommandHandler,
	title, description, difficulty string,
	reward, durationMinutes int,
	creator string,
	targetLocation, executionLocation kernel.GeoCoordinate,
	equipment, skills []string,
) (quest.Quest, error) {
	cmd := commands.CreateQuestCommand{
		Title:             title,
		Description:       description,
		Difficulty:        difficulty,
		Reward:            reward,
		DurationMinutes:   durationMinutes,
		TargetLocation:    targetLocation,
		ExecutionLocation: executionLocation,
		Creator:           creator,
		Equipment:         equipment,
		Skills:            skills,
	}

	return handler.Handle(ctx, cmd)
}

// AssignQuestStep выполняет назначение квеста пользователю
func AssignQuestStep(
	ctx context.Context,
	handler commands.AssignQuestCommandHandler,
	questID uuid.UUID,
	userID string,
) (commands.AssignQuestResult, error) {
	cmd := commands.AssignQuestCommand{
		ID:     questID,
		UserID: userID,
	}

	return handler.Handle(ctx, cmd)
}

// ChangeQuestStatusStep изменяет статус квеста
func ChangeQuestStatusStep(
	ctx context.Context,
	handler commands.ChangeQuestStatusCommandHandler,
	questRepo ports.QuestRepository,
	questID uuid.UUID,
	newStatus quest.Status,
) (quest.Quest, error) {
	cmd := commands.ChangeQuestStatusCommand{
		QuestID: questID,
		Status:  newStatus,
	}

	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		return quest.Quest{}, err
	}

	// Для тестов получаем полный квест
	return questRepo.GetByID(ctx, result.ID)
}

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
