package cases

import (
	"quest-manager/cmd"

	"github.com/stretchr/testify/suite"
)

// DefaultSuite basic test suite for integration tests
type DefaultSuite struct {
	SuiteDIContainer
	TestDIContainer
}

// NewDefault creates new DefaultSuite
func NewDefault(s suite.TestingSuite) DefaultSuite {
	return DefaultSuite{
		SuiteDIContainer: NewSuite(s),
	}
}

// SetupSuite initializes resources before running all tests in the suite
func (s *DefaultSuite) SetupSuite() {
	s.TestDIContainer = NewTestDIContainer(s.SuiteDIContainer)

	// Run migrations
	cmd.MustAutoMigrate(s.TestDIContainer.DB)
}

// TearDownSuite cleans up resources after completing all tests in the suite
func (s *DefaultSuite) TearDownSuite() {
	s.TestDIContainer.TearDownTest()
}

// SetupTest prepares state before each test
func (s *DefaultSuite) SetupTest() {
	// Clean database before each test
	err := s.TestDIContainer.CleanupDatabase()
	s.Require().NoError(err, "Failed to cleanup database")

	// Recreate TestDIContainer for each test to avoid transaction issues
	s.TestDIContainer = NewTestDIContainer(s.SuiteDIContainer)
}

// TearDownTest cleans state after each test
func (s *DefaultSuite) TearDownTest() {
	// Wait for event processing completion
	s.TestDIContainer.WaitForEventProcessing()
}
