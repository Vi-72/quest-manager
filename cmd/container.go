package cmd

import (
	"context"
	"fmt"
	"gorm.io/gorm"

	v1 "quest-manager/api/http/quests/v1"
	httphandlers "quest-manager/internal/adapters/in/http"
	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"
)

// Container holds all application dependencies and provides access to services.
type Container struct {
	configs        Config
	db             *gorm.DB
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
	authClient     ports.AuthClient
	closers        []Closer
}

// NewContainer creates a new dependency injection container.
// All heavy initialization is deferred to Build() method.
func NewContainer(configs Config, db *gorm.DB) (*Container, error) {
	unitOfWork, err := postgres.NewUnitOfWork(db)
	if err != nil {
		return nil, fmt.Errorf("create unit of work: %w", err)
	}

	eventPublisher, err := eventrepo.NewRepository(
		unitOfWork.(ports.Tracker),
		configs.EventGoroutineLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("create event publisher: %w", err)
	}

	return &Container{
		configs:        configs,
		db:             db,
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
	}, nil
}

// Cfg returns configuration.
func (c *Container) Cfg() Config { return c.configs }

// DB returns database connection.
func (c *Container) DB() *gorm.DB { return c.db }

// GetUnitOfWork returns the single UnitOfWork instance.
func (c *Container) GetUnitOfWork() ports.UnitOfWork { return c.unitOfWork }

// EventPublisher returns EventPublisher.
func (c *Container) EventPublisher() ports.EventPublisher { return c.eventPublisher }

// GetAuthClient returns auth client (lazy-initialized via factory).
func (c *Container) GetAuthClient(ctx context.Context) ports.AuthClient {
	if c.authClient != nil {
		return c.authClient
	}
	client, conn, _ := c.configs.AuthFactory.Create(ctx)
	c.RegisterCloser(CloserFunc(conn.Close))

	c.authClient = client
	return c.authClient
}

// QuestRepository returns repository from the single UoW.
func (c *Container) QuestRepository() ports.QuestRepository {
	return c.unitOfWork.QuestRepository()
}

// LocationRepository returns repository from the single UoW.
func (c *Container) LocationRepository() ports.LocationRepository {
	return c.unitOfWork.LocationRepository()
}

// Handlers groups all command/query handlers for API wiring.
type Handlers struct {
	CreateQuest       commands.CreateQuestCommandHandler
	ListQuests        queries.ListQuestsQueryHandler
	GetQuestByID      queries.GetQuestByIDQueryHandler
	ChangeQuestStatus commands.ChangeQuestStatusCommandHandler
	AssignQuest       commands.AssignQuestCommandHandler
	SearchByRadius    queries.SearchQuestsByRadiusQueryHandler
	ListAssigned      queries.ListAssignedQuestsQueryHandler
}

// Handlers initializes all application handlers.
func (c *Container) Handlers() Handlers {
	return Handlers{
		CreateQuest:       commands.NewCreateQuestCommandHandler(c.unitOfWork, c.eventPublisher),
		ListQuests:        queries.NewListQuestsQueryHandler(c.QuestRepository()),
		GetQuestByID:      queries.NewGetQuestByIDQueryHandler(c.QuestRepository()),
		ChangeQuestStatus: commands.NewChangeQuestStatusCommandHandler(c.unitOfWork, c.eventPublisher),
		AssignQuest:       commands.NewAssignQuestCommandHandler(c.unitOfWork, c.eventPublisher),
		SearchByRadius:    queries.NewSearchQuestsByRadiusQueryHandler(c.QuestRepository()),
		ListAssigned:      queries.NewListAssignedQuestsQueryHandler(c.QuestRepository()),
	}
}

// NewApiHandler aggregates all HTTP handlers.
func (c *Container) NewApiHandler() (v1.StrictServerInterface, error) {
	h := c.Handlers()
	return httphandlers.NewApiHandler(
		h.CreateQuest,
		h.ListQuests,
		h.GetQuestByID,
		h.ChangeQuestStatus,
		h.SearchByRadius,
		h.ListAssigned,
		h.AssignQuest,
	)
}

// --- utility ---

// CloserFunc allows inline closing functions to implement io.Closer.
type CloserFunc func() error

func (f CloserFunc) Close() error { return f() }
