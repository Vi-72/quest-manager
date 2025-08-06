package quest_e2e_tests

import (
	"fmt"
	"testing"

	"quest-manager/tests/integration/tests"

	"github.com/stretchr/testify/suite"
)

// E2ESuite provides the base setup for all E2E tests
type E2ESuite struct {
	suite.Suite
	tests.DefaultSuite
}

func (s *E2ESuite) SetupSuite() {
	s.DefaultSuite = tests.NewDefault(s)
	s.DefaultSuite.SetupSuite()
}

func (s *E2ESuite) SetupTest() {
	s.DefaultSuite.SetupTest()
}

func (s *E2ESuite) TearDownTest() {
	s.DefaultSuite.TearDownTest()
}

func (s *E2ESuite) TearDownSuite() {
	s.DefaultSuite.TearDownSuite()
}

// TestE2EOperations is a placeholder test runner
func TestE2EOperations(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}

// Example of how to run a simple smoke test
func (s *E2ESuite) TestE2ESmokeTest() {
	// Simple test to verify the test environment is working
	s.Assert().NotNil(s.TestDIContainer, "DI Container should be initialized")
	s.Assert().NotNil(s.TestDIContainer.HTTPRouter, "HTTP Router should be initialized")
	s.Assert().NotNil(s.TestDIContainer.QuestRepository, "Quest Repository should be initialized")
	s.Assert().NotNil(s.TestDIContainer.LocationRepository, "Location Repository should be initialized")

	fmt.Println("🎯 E2E Test Environment is ready!")
}
