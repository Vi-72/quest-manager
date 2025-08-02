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
	container := NewTestDIContainer(s.SuiteDIContainer, false)
	cmd.MustAutoMigrate(container.DB)
	container.TearDownTest()
}

// TearDownSuite cleans up resources after completing all tests in the suite
func (s *DefaultSuite) TearDownSuite() {}

// SetupTest prepares state before each test
func (s *DefaultSuite) SetupTest() {
	s.TestDIContainer = NewTestDIContainer(s.SuiteDIContainer, true)
}

// TearDownTest cleans state after each test
func (s *DefaultSuite) TearDownTest() {
	s.TestDIContainer.WaitForEventProcessing(0)
	s.TestDIContainer.TearDownTest()
}
