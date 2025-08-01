package quest_operations

import (
	"testing"

	"quest-manager/tests/integration/cases"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	cases.DefaultSuite
}

func TestQuestOperations(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	s.DefaultSuite = cases.NewDefault(s)
	s.DefaultSuite.SetupSuite()
}

func (s *Suite) SetupTest() {
	s.DefaultSuite.SetupTest()
}

func (s *Suite) TearDownTest() {
	s.DefaultSuite.TearDownTest()
}

func (s *Suite) TearDownSuite() {
	s.DefaultSuite.TearDownSuite()
}
