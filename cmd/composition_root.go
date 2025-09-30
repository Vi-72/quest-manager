package cmd

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"

	authv1 "github.com/Vi-72/quest-auth/api/grpc/sdk/go/auth/v1"

	v1 "quest-manager/api/http/quests/v1"
	httphandlers "quest-manager/internal/adapters/in/http"
	authclient "quest-manager/internal/adapters/out/client/auth"
	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"
)

type CompositionRoot struct {
	configs        Config
	db             *gorm.DB
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
	closers        []Closer

	// auth
	authConn      *grpc.ClientConn
	authSDKClient authv1.AuthServiceClient
	authClient    ports.AuthClient
}

func NewCompositionRoot(configs Config, db *gorm.DB) *CompositionRoot {
	// Create Unit of Work once during initialization
	unitOfWork, err := postgres.NewUnitOfWork(db)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}

	// Create EventPublisher with same Tracker as UoW for transactionality
	eventPublisher, err := eventrepo.NewRepository(unitOfWork.(ports.Tracker), configs.EventGoroutineLimit)
	if err != nil {
		log.Fatalf("cannot create EventPublisher: %v", err)
	}

	cr := &CompositionRoot{
		configs:        configs,
		db:             db,
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
	}

	// ---- wire Auth gRPC client (optional: if AUTH_GRPC provided)
	if addr := configs.AuthGRPC; addr != "" {
		dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Можно заменить на mTLS/creds при необходимости
		conn, err := grpc.DialContext(
			dialCtx,
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			log.Fatalf("failed to dial auth gRPC at %s: %v", addr, err)
		}
		cr.authConn = conn
		cr.RegisterCloser(connCloser{conn}) // закроем при shutdown

		cr.authSDKClient = authv1.NewAuthServiceClient(conn)
		cr.authClient = authclient.NewUserAuthClient(cr.authSDKClient)
	}

	return cr
}

// helper to close grpc.Conn via our Closer stack
type connCloser struct{ *grpc.ClientConn }

func (c connCloser) Close() error { return c.ClientConn.Close() }

// GetUnitOfWork returns the single UnitOfWork instance
func (cr *CompositionRoot) GetUnitOfWork() ports.UnitOfWork {
	return cr.unitOfWork
}

// EventPublisher returns EventPublisher
func (cr *CompositionRoot) EventPublisher() ports.EventPublisher {
	return cr.eventPublisher
}

// QuestRepository returns repository from the single UoW
func (cr *CompositionRoot) QuestRepository() ports.QuestRepository {
	return cr.unitOfWork.QuestRepository()
}

// LocationRepository returns repository from the single UoW
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
func (cr *CompositionRoot) NewApiHandler() v1.StrictServerInterface {
	handlers, err := httphandlers.NewApiHandler(
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
