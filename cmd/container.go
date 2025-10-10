package cmd

import (
	"fmt"

	authv1 "github.com/Vi-72/quest-auth/api/grpc/sdk/go/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gorm.io/gorm"

	v1 "quest-manager/api/http/quests/v1"
	httphandlers "quest-manager/internal/adapters/in/http"
	authclient "quest-manager/internal/adapters/out/client/auth"
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
// Initializes all dependencies including auth client (eager initialization).
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

	container := &Container{
		configs:        configs,
		db:             db,
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
	}

	if !configs.Middleware.DevAuth.Enabled {
		authClient, _ := container.createAuthClient()
		container.authClient = authClient
	}

	return container, nil
}

// Cfg returns configuration.
func (c *Container) Cfg() Config { return c.configs }

// DB returns database connection.
func (c *Container) DB() *gorm.DB { return c.db }

// GetUnitOfWork returns the single UnitOfWork instance.
func (c *Container) GetUnitOfWork() ports.UnitOfWork { return c.unitOfWork }

// EventPublisher returns EventPublisher.
func (c *Container) EventPublisher() ports.EventPublisher { return c.eventPublisher }

// GetAuthClient returns auth client (initialized in NewContainer or injected via SetAuthClient).
func (c *Container) GetAuthClient() ports.AuthClient {
	return c.authClient
}

// SetAuthClient allows injecting a custom auth client (for testing).
func (c *Container) SetAuthClient(client ports.AuthClient) {
	c.authClient = client
}

// createAuthClient creates and initializes auth gRPC client (internal helper).
func (c *Container) createAuthClient() (ports.AuthClient, error) {
	// Create gRPC connection
	conn, err := c.createGRPCConnection(c.configs.AuthGRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth gRPC connection: %w", err)
	}

	// Register connection for cleanup
	c.RegisterCloser(CloserFunc(conn.Close))

	// Create auth gRPC client and wrap it
	grpcClient := authv1.NewAuthServiceClient(conn)
	authClient := authclient.NewUserAuthClient(grpcClient)

	return authClient, nil
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

// createGRPCConnection creates a gRPC client connection with insecure credentials.
func (c *Container) createGRPCConnection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// CloserFunc allows inline closing functions to implement io.Closer.
type CloserFunc func() error

func (f CloserFunc) Close() error { return f() }
