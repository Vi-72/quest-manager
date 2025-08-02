package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/adapters/out/postgres/locationrepo"
	"quest-manager/internal/adapters/out/postgres/questrepo"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
)

// RepositoryTestSuite provides shared setup for repository tests
type RepositoryTestSuite struct {
	suite.Suite
	db             *gorm.DB
	tracker        ports.Tracker
	questRepo      ports.QuestRepository
	locationRepo   ports.LocationRepository
	eventPublisher ports.EventPublisher
}

func (suite *RepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite for fast tests
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Quiet logs for tests
	})
	suite.Require().NoError(err)
	suite.db = db

	// Auto-migrate all tables
	err = suite.migrateSchema()
	suite.Require().NoError(err)

	// Create tracker (unit of work)
	tracker := NewTestTracker(db)
	suite.tracker = tracker

	// Create repositories
	questRepo, err := questrepo.NewRepository(tracker)
	suite.Require().NoError(err)
	suite.questRepo = questRepo

	locationRepo, err := locationrepo.NewRepository(tracker)
	suite.Require().NoError(err)
	suite.locationRepo = locationRepo

	// Mock event publisher for repository tests
	suite.eventPublisher = &NoOpEventPublisher{}
}

func (suite *RepositoryTestSuite) SetupTest() {
	// Clean database before each test
	suite.cleanDatabase()
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

func (suite *RepositoryTestSuite) migrateSchema() error {
	// Import DTO structs and migrate
	return suite.db.AutoMigrate(
		&questrepo.QuestDTO{},
		&locationrepo.LocationDTO{},
		&eventrepo.EventDTO{},
	)
}

func (suite *RepositoryTestSuite) cleanDatabase() {
	// Clean all tables in reverse order (respect foreign keys)
	suite.db.Exec("DELETE FROM quests")
	suite.db.Exec("DELETE FROM locations")
	suite.db.Exec("DELETE FROM events")
}

// TestTracker implements ports.Tracker for testing
type TestTracker struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewTestTracker(db *gorm.DB) *TestTracker {
	return &TestTracker{db: db}
}

func (t *TestTracker) Begin(ctx context.Context) error {
	t.tx = t.db.Begin()
	return t.tx.Error
}

func (t *TestTracker) Commit(ctx context.Context) error {
	if t.tx == nil {
		return nil
	}
	err := t.tx.Commit().Error
	t.tx = nil
	return err
}

func (t *TestTracker) Rollback() error {
	if t.tx == nil {
		return nil
	}
	err := t.tx.Rollback().Error
	t.tx = nil
	return err
}

func (t *TestTracker) Tx() *gorm.DB {
	if t.tx != nil {
		return t.tx
	}
	return t.db
}

func (t *TestTracker) InTx() bool {
	return t.tx != nil
}

func (t *TestTracker) Db() *gorm.DB {
	return t.db
}

// NoOpEventPublisher is a mock event publisher for repository tests
type NoOpEventPublisher struct{}

func (n *NoOpEventPublisher) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	// Do nothing - repository tests don't need event publishing
	return nil
}

func (n *NoOpEventPublisher) PublishAsync(ctx context.Context, events ...ddd.DomainEvent) {
	// Do nothing - repository tests don't need event publishing
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
