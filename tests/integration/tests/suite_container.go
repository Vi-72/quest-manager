package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// SuiteDIContainer предоставляет базовую функциональность для тестовых наборов
type SuiteDIContainer struct {
	testingSuite suite.TestingSuite
}

// T возвращает тестинг объект
func (s SuiteDIContainer) T() *testing.T {
	return s.testingSuite.T()
}

// Require возвращает require объект для assertions
func (s SuiteDIContainer) Require() *require.Assertions {
	return require.New(s.testingSuite.T())
}

// Assert возвращает assert объект для assertions
func (s SuiteDIContainer) Assert() *assert.Assertions {
	return assert.New(s.testingSuite.T())
}

// NewSuite создает новый контейнер для тестового набора
func NewSuite(s suite.TestingSuite) SuiteDIContainer {
	return SuiteDIContainer{
		testingSuite: s,
	}
}
