# Integration Tests

Эта папка содержит интеграционные тесты для Quest Manager Service.
oapi-codegen -config configs/server.cfg.yaml api/openapi/openapi.yml

## Структура

```
tests/integration/
├── cases/                     # Тестовые наборы
│   ├── quest_operations/      # Тесты операций с квестами
│   ├── test_container.go      # DI контейнер для тестов
│   ├── suite_container.go     # Базовый контейнер для test suites
│   └── default.go            # Базовый тестовый набор
├── core/                      # Основные компоненты
│   ├── assertions/           # Пользовательские assertions
│   ├── case_steps/          # Переиспользуемые шаги тестирования
│   ├── storage/             # Утилиты для работы с базой данных
│   └── test_data_generators/ # Генераторы тестовых данных
├── mock/                     # Mock объекты (если нужны)
└── README.md                # Этот файл
```

## Компоненты

### Test Container (`test_container.go`)
Центральный DI контейнер, который:
- Инициализирует тестовую базу данных
- Создает все необходимые репозитории и use cases
- Управляет жизненным циклом ресурсов
- Обеспечивает изоляцию тестов

### Case Steps (`core/case_steps/`)
Переиспользуемые шаги для тестирования:
- `quest_operations.go` - операции с квестами (создание, назначение, изменение статуса)
- `http_requests.go` - HTTP запросы к API

### Storage (`core/storage/`)
Утилиты для прямого доступа к базе данных в тестах:
- `quest_storage.go` - работа с квестами
- `location_storage.go` - работа с локациями  
- `event_storage.go` - работа с событиями

### Assertions (`core/assertions/`)
Пользовательские проверки:
- `quest_assertions.go` - проверки состояния квестов и событий

### Test Data Generators (`core/test_data_generators/`)
Генераторы тестовых данных:
- `quest_generator.go` - создание тестовых данных для квестов

## Как запускать тесты

### Подготовка

1. Убедитесь что PostgreSQL запущен:
```bash
docker compose up -d postgres
```

2. Создайте тестовую базу данных:
```sql
CREATE DATABASE quest_manager_test;
```

### Запуск тестов

Запуск всех интеграционных тестов:
```bash
go test -tags=integration ./tests/integration/...
```

Запуск конкретного тестового набора:
```bash
go test -tags=integration ./tests/integration/cases/quest/
```

Запуск с подробным выводом:
```bash
go test -tags=integration -v ./tests/integration/...
```

## Примеры тестов

### Test Lifecycle (`quest_lifecycle_test.go`)
Тестирует полный жизненный цикл квеста:
- Создание квеста
- Назначение пользователю
- Изменение статуса
- Проверка событий

### API Tests (`quest_api_test.go`)  
Тестирует HTTP API:
- Создание квеста через POST /api/v1/quests
- Назначение через POST /api/v1/quests/{id}/assign
- Изменение статуса через PATCH /api/v1/quests/{id}/status
- Получение квеста через GET /api/v1/quests/{id}

## Конфигурация

Тесты используют отдельную конфигурацию:
- База данных: `quest_manager_test`
- Порт: `8081` (вместо 8080)
- Лимит горутин событий: `3` (вместо 5)

## Принципы

1. **Изоляция**: Каждый тест очищает базу данных
2. **Детерминизм**: Используются фиксированные данные где возможно
3. **Быстрота**: Тесты должны выполняться быстро
4. **Читаемость**: Тесты должны быть понятными и хорошо структурированными
5. **Покрытие**: Тестируем важные пути и граничные случаи

## Добавление новых тестов

1. Создайте новый файл в `cases/` с суффиксом `_test.go`
2. Добавьте build tag `//go:build integration`
3. Используйте `DefaultSuite` как базу
4. Создавайте переиспользуемые шаги в `case_steps/`
5. Добавляйте генераторы данных в `test_data_generators/`
6. Используйте существующие assertions или создавайте новые

Пример структуры нового теста:
```go
//go:build integration

package newfeature

import (
    "testing"
    "github.com/stretchr/testify/suite"
    "quest-manager/tests/integration/cases"
)

type NewFeatureSuite struct {
    suite.Suite
    cases.DefaultSuite
}

func TestNewFeature(t *testing.T) {
    suite.Run(t, new(NewFeatureSuite))
}

func (s *NewFeatureSuite) SetupSuite() {
    s.DefaultSuite.SetupSuite()
}

func (s *NewFeatureSuite) TestSomething() {
    // Arrange
    // Act  
    // Assert
}
```