package tests

import (
	"quest-manager/cmd"

	"github.com/stretchr/testify/suite"
)

// DefaultSuite basic test suite for integration tests
type DefaultSuite struct {
	SuiteDIContainer
	TestDIContainer *TestContainer
}

// NewDefault creates new DefaultSuite
func NewDefault(s suite.TestingSuite) DefaultSuite {
	return DefaultSuite{
		SuiteDIContainer: NewSuite(s),
	}
}

// SetupSuite initializes resources before running all tests in the suite
func (s *DefaultSuite) SetupSuite() {
	// Get DB from the test environment
	// Use test database URL - in real tests this would come from environment
	testDBURL := "postgres://postgres:password@localhost:5432/quest_manager_test?sslmode=disable"
	db, _, err := cmd.MustConnectDB(testDBURL)
	s.Require().NoError(err, "Failed to connect to test database")

	s.TestDIContainer = NewTestContainer(db)

	// Run migrations
	cmd.MustAutoMigrate(s.TestDIContainer.DB)
}

// TearDownSuite cleans up resources after completing all tests in the suite
func (s *DefaultSuite) TearDownSuite() {
	s.TestDIContainer.CleanupAll()
}

// SetupTest prepares state before each test
func (s *DefaultSuite) SetupTest() {
	// Clean database before each test
	s.TestDIContainer.CleanupAll()
}

// TearDownTest cleans state after each test
func (s *DefaultSuite) TearDownTest() {
	// Wait for event processing completion
	s.TestDIContainer.WaitForEventProcessing(0)
}
