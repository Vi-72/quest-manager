package cmd

import (
	"log"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http"
	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs           Config
	db                *gorm.DB
	unitOfWorkFactory ports.UnitOfWorkFactory
	closers           []Closer
}

func NewCompositionRoot(configs Config, db *gorm.DB) *CompositionRoot {
	factory := func() (ports.UnitOfWork, ports.EventPublisher, error) {
		unitOfWork, err := postgres.NewUnitOfWork(db)
		if err != nil {
			return nil, nil, err
		}

		tracker, ok := unitOfWork.(ports.Tracker)
		if !ok {
			return nil, nil, errs.WrapInfrastructureError("unit of work does not implement tracker", nil)
		}

		eventPublisher, err := eventrepo.NewRepository(tracker, configs.EventGoroutineLimit)
		if err != nil {
			return nil, nil, err
		}

		return unitOfWork, eventPublisher, nil
	}

	if _, _, err := factory(); err != nil {
		log.Fatalf("cannot initialize application dependencies: %v", err)
	}

	return &CompositionRoot{
		configs:           configs,
		db:                db,
		unitOfWorkFactory: factory,
	}
}

// UnitOfWorkFactory returns a factory that produces fresh UnitOfWork instances
// together with their request-scoped event publishers.
func (cr *CompositionRoot) UnitOfWorkFactory() ports.UnitOfWorkFactory {
	return cr.unitOfWorkFactory
}

// NewCreateQuestCommandHandler creates a handler for creating quests.
func (cr *CompositionRoot) NewCreateQuestCommandHandler() commands.CreateQuestCommandHandler {
	return commands.NewCreateQuestCommandHandler(cr.UnitOfWorkFactory())
}

// NewListQuestsQueryHandler creates a handler for listing quests.
func (cr *CompositionRoot) NewListQuestsQueryHandler() queries.ListQuestsQueryHandler {
	return queries.NewListQuestsQueryHandler(cr.UnitOfWorkFactory())
}

// NewGetQuestByIDQueryHandler creates a handler for fetching a quest by its ID.
func (cr *CompositionRoot) NewGetQuestByIDQueryHandler() queries.GetQuestByIDQueryHandler {
	return queries.NewGetQuestByIDQueryHandler(cr.UnitOfWorkFactory())
}

// NewChangeQuestStatusHandler creates a handler for changing quest status.
func (cr *CompositionRoot) NewChangeQuestStatusHandler() commands.ChangeQuestStatusCommandHandler {
	return commands.NewChangeQuestStatusCommandHandler(cr.UnitOfWorkFactory())
}

// NewAssignQuestCommandHandler creates a handler for assigning a quest.
func (cr *CompositionRoot) NewAssignQuestCommandHandler() commands.AssignQuestCommandHandler {
	return commands.NewAssignQuestCommandHandler(cr.UnitOfWorkFactory())
}

// NewSearchQuestsByRadiusQueryHandler creates a handler for searching quests in a radius.
func (cr *CompositionRoot) NewSearchQuestsByRadiusQueryHandler() queries.SearchQuestsByRadiusQueryHandler {
	return queries.NewSearchQuestsByRadiusQueryHandler(cr.UnitOfWorkFactory())
}

// NewListAssignedQuestsQueryHandler creates a handler for listing quests assigned to a user.
func (cr *CompositionRoot) NewListAssignedQuestsQueryHandler() queries.ListAssignedQuestsQueryHandler {
	return queries.NewListAssignedQuestsQueryHandler(cr.UnitOfWorkFactory())
}

// NewApiHandler aggregates all HTTP handlers.
func (cr *CompositionRoot) NewApiHandler() v1.StrictServerInterface {
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
