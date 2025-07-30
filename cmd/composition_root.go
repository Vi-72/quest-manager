package cmd

import (
	"log"
	"quest-manager/internal/adapters/in/http"
	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/generated/servers"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs Config
	db      *gorm.DB

	closers []Closer
}

func NewCompositionRoot(configs Config, db *gorm.DB) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
		db:      db,
	}
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(cr.db)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}

func (cr *CompositionRoot) QuestRepository() ports.QuestRepository {
	return cr.NewUnitOfWork().QuestRepository()
}

// NewCreateQuestCommandHandler creates a handler for creating quests.
func (cr *CompositionRoot) NewCreateQuestCommandHandler() commands.CreateQuestCommandHandler {
	return commands.NewCreateQuestCommandHandler(cr.QuestRepository())
}

// NewListQuestsQueryHandler creates a handler for listing quests.
func (cr *CompositionRoot) NewListQuestsQueryHandler() queries.ListQuestsQueryHandler {
	return queries.NewListQuestsQueryHandler(cr.QuestRepository())
}

// NewGetQuestByIDQueryHandler creates a handler for fetching a quest by its ID.
func (cr *CompositionRoot) NewGetQuestByIDQueryHandler() queries.GetQuestByIDQueryHandler {
	return queries.NewGetQuestByIDQueryHandler(cr.QuestRepository())
}

// NewChangeQuestStatusCommandHandler creates a handler for changing quest status.
func (cr *CompositionRoot) NewChangeQuestStatusCommandHandler() commands.ChangeQuestStatusCommandHandler {
	return commands.NewChangeQuestStatusCommandHandler(cr.QuestRepository())
}

// NewAssignQuestCommandHandler creates a handler for assigning a quest.
func (cr *CompositionRoot) NewAssignQuestCommandHandler() commands.AssignQuestCommandHandler {
	return commands.NewAssignQuestCommandHandler(cr.QuestRepository())
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
		cr.NewChangeQuestStatusCommandHandler(),
		cr.NewSearchQuestsByRadiusQueryHandler(),
		cr.NewListAssignedQuestsQueryHandler(),
		cr.NewAssignQuestCommandHandler(),
	)
	if err != nil {
		log.Fatalf("Error initializing HTTP Server: %v", err)
	}
	return handlers
}
