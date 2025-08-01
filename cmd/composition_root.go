package cmd

import (
	"log"
	"quest-manager/internal/adapters/in/http"
	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/generated/servers"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs        Config
	db             *gorm.DB
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
	closers        []Closer
}

func NewCompositionRoot(configs Config, db *gorm.DB) *CompositionRoot {
	// Создаем Unit of Work один раз при инициализации
	unitOfWork, err := postgres.NewUnitOfWork(db)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}

	// Создаем EventPublisher с тем же Tracker что и UoW для транзакционности
	eventPublisher, err := eventrepo.NewRepository(unitOfWork.(ports.Tracker), configs.EventGoroutineLimit)
	if err != nil {
		log.Fatalf("cannot create EventPublisher: %v", err)
	}

	return &CompositionRoot{
		configs:        configs,
		db:             db,
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
	}
}

// GetUnitOfWork возвращает единственный экземпляр UnitOfWork
func (cr *CompositionRoot) GetUnitOfWork() ports.UnitOfWork {
	return cr.unitOfWork
}

// EventPublisher возвращает EventPublisher
func (cr *CompositionRoot) EventPublisher() ports.EventPublisher {
	return cr.eventPublisher
}

// QuestRepository возвращает репозиторий из единственного UoW
func (cr *CompositionRoot) QuestRepository() ports.QuestRepository {
	return cr.unitOfWork.QuestRepository()
}

// LocationRepository возвращает репозиторий из единственного UoW
func (cr *CompositionRoot) LocationRepository() ports.LocationRepository {
	return cr.unitOfWork.LocationRepository()
}

// NewCreateQuestCommandHandler creates a handler for creating quests.
func (cr *CompositionRoot) NewCreateQuestCommandHandler() commands.CreateQuestCommandHandler {
	return commands.NewCreateQuestCommandHandler(cr.GetUnitOfWork(), cr.EventPublisher())
}

// NewListQuestsQueryHandler creates a handler for listing quests.
func (cr *CompositionRoot) NewListQuestsQueryHandler() queries.ListQuestsQueryHandler {
	return queries.NewListQuestsQueryHandler(cr.QuestRepository())
}

// NewGetQuestByIDQueryHandler creates a handler for fetching a quest by its ID.
func (cr *CompositionRoot) NewGetQuestByIDQueryHandler() queries.GetQuestByIDQueryHandler {
	return queries.NewGetQuestByIDQueryHandler(cr.QuestRepository())
}

// NewChangeQuestStatusHandler creates a handler for changing quest status.
func (cr *CompositionRoot) NewChangeQuestStatusHandler() commands.ChangeQuestStatusCommandHandler {
	return commands.NewChangeQuestStatusCommandHandler(cr.GetUnitOfWork(), cr.EventPublisher())
}

// NewAssignQuestCommandHandler creates a handler for assigning a quest.
func (cr *CompositionRoot) NewAssignQuestCommandHandler() commands.AssignQuestCommandHandler {
	return commands.NewAssignQuestCommandHandler(cr.GetUnitOfWork(), cr.EventPublisher())
}

// NewSearchQuestsByRadiusQueryHandler creates a handler for searching quests in a radius.
func (cr *CompositionRoot) NewSearchQuestsByRadiusQueryHandler() queries.SearchQuestsByRadiusQueryHandler {
	return queries.NewSearchQuestsByRadiusQueryHandler(cr.QuestRepository())
}

// NewListAssignedQuestsQueryHandler creates a handler for listing quests assigned to a user.
func (cr *CompositionRoot) NewListAssignedQuestsQueryHandler() queries.ListAssignedQuestsQueryHandler {
	return queries.NewListAssignedQuestsQueryHandler(cr.QuestRepository())
}

// NewApiHandler aggregates all HTTP handlers.
func (cr *CompositionRoot) NewApiHandler() servers.StrictServerInterface {
	handlers, err := http.NewApiHandler(
		cr.NewCreateQuestCommandHandler(),
		cr.NewListQuestsQueryHandler(),
		cr.NewGetQuestByIDQueryHandler(),
		cr.NewChangeQuestStatusHandler(),
		cr.NewSearchQuestsByRadiusQueryHandler(),
		cr.NewListAssignedQuestsQueryHandler(),
		cr.NewAssignQuestCommandHandler(),
	)
	if err != nil {
		log.Fatalf("Error initializing HTTP Server: %v", err)
	}
	return handlers
}
