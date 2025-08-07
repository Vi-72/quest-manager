# Integration Tests

Эта папка содержит интеграционные тесты для Quest Manager Service.
oapi-codegen -config configs/server.cfg.yaml api/openapi/openapi.yml

## Структура

```
tests/integration/
├── tests/                     # 🧪 Группы интеграционных тестов
│   ├── quest_e2e_tests/       # E2E тесты полного цикла квестов
│   ├── quest_http_tests/      # HTTP API layer тесты
│   ├── quest_handler_tests/   # Application handler тесты
│   ├── repository_tests/      # Infrastructure repository тесты
│   ├── test_container.go      # DI контейнер для тестов
│   ├── suite_container.go     # Базовый контейнер для test suites
│   └── default.go            # Базовый тестовый набор
├── core/                      # 🔧 Переиспользуемые компоненты
│   ├── assertions/           # Пользовательские проверки
│   │   ├── quest_assign_assertions.go
│   │   ├── quest_e2e_assertions.go
│   │   ├── quest_field_assertions.go
│   │   ├── quest_handler_assertions.go    # ✨ Новый
│   │   ├── quest_http_assertions.go
│   │   ├── quest_list_assertions.go
│   │   └── quest_single_assertions.go
│   ├── case_steps/          # Переиспользуемые шаги тестирования
│   │   ├── database_steps.go
│   │   ├── http_requests.go
│   │   ├── quest_creation.go
│   │   ├── quest_queries.go
│   │   └── quest_status.go
│   ├── storage/             # Прямой доступ к БД
│   │   └── event_storage.go
│   └── test_data_generators/ # Генераторы тестовых данных
│       └── quest_generator.go
└── README.md                # Этот файл
```

## Компоненты

### Test Container (`test_container.go`)
Центральный DI контейнер, который:
- Инициализирует тестовую базу данных
- Создает все необходимые репозитории и use cases
- Управляет жизненным циклом ресурсов
- Обеспечивает изоляцию тестов

### Test Groups (`tests/`)
Организованы по слоям архитектуры:

#### **🌐 E2E Tests** (`quest_e2e_tests/`)
Тесты полного цикла квеста от создания до завершения:
- Создание квеста через Handler, назначение через API
- Проверка событий и состояния в БД
- Smoke тесты среды

#### **🌍 HTTP Tests** (`quest_http_tests/`)
Тесты HTTP API слоя:
- Валидация входных данных
- Коды ответов и форматы JSON
- Error handling и edge cases

#### **⚙️ Handler Tests** (`quest_handler_tests/`)
Тесты Application слоя (use cases):
- Бизнес-логика без HTTP слоя
- Оркестрация команд и queries
- Domain events генерация

#### **🗄️ Repository Tests** (`repository_tests/`)
Тесты Infrastructure слоя:
- PostgreSQL интеграция
- CRUD операции
- Transaction handling

### Case Steps (`core/case_steps/`)
Переиспользуемые шаги для всех типов тестов:
- `quest_creation.go` - создание квестов
- `quest_queries.go` - получение данных квестов
- `quest_status.go` - операции изменения статуса
- `http_requests.go` - HTTP запросы к API
- `database_steps.go` - прямая работа с БД

### Storage (`core/storage/`)
Утилиты для прямого доступа к базе данных в тестах:
- `event_storage.go` - работа с событиями

### Assertions (`core/assertions/`)
Специализированные проверки по слоям:
- `quest_e2e_assertions.go` - E2E сценарии
- `quest_http_assertions.go` - HTTP responses  
- `quest_handler_assertions.go` - Handler логика ✨
- `quest_assign_assertions.go` - Назначение квестов
- `quest_field_assertions.go` - Поля и валидация
- `quest_list_assertions.go` - Списки квестов
- `quest_single_assertions.go` - Отдельные квесты

### Test Data Generators (`core/test_data_generators/`)
Генераторы тестовых данных:
- `quest_generator.go` - создание тестовых данных для квестов всех типов

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

#### Все интеграционные тесты:
```bash
make test-integration
# или
go test -tags=integration ./tests/integration/... -v
```

#### По группам тестов:
```bash
# E2E тесты
go test -tags=integration ./tests/integration/tests/quest_e2e_tests -v

# HTTP API тесты  
go test -tags=integration ./tests/integration/tests/quest_http_tests -v

# Handler тесты
go test -tags=integration ./tests/integration/tests/quest_handler_tests -v

# Repository тесты
make test-repository
# или
go test -tags=integration ./tests/integration/tests/repository_tests -v
```

#### С анализом покрытия:
```bash
make test-coverage-integration
```

## Примеры тестов

### E2E Tests (`quest_e2e_tests/`)
```go
// assign_quest_e2e_test.go
func (s *Suite) TestCreateThroughHandlerAssignThroughAPI() {
    // Создание через Handler слой
    createdQuest := casesteps.CreateRandomQuestStep(...)
    
    // Назначение через HTTP API
    response := casesteps.AssignQuestHTTPStep(...)
    
    // Проверка в БД и событий
    assertions.VerifyQuestAssignedCorrectly(...)
}
```

### HTTP API Tests (`quest_http_tests/`)
```go  
// create_quest_http_test.go
func (s *Suite) TestCreateQuestHTTP() {
    // Подготовка HTTP запроса
    requestData := testdatagenerators.ValidHTTPQuestData()
    
    // HTTP POST /api/v1/quests
    response := casesteps.CreateQuestHTTPStep(...)
    
    // Проверка HTTP response
    assertions.VerifyHTTPCreateResponse(...)
}
```

### Handler Tests (`quest_handler_tests/`)
```go
// create_quest_test.go  
func (s *Suite) TestCreateQuestWithAllParameters() {
    // Подготовка команды
    questData := testdatagenerators.SimpleQuestData(...)
    
    // Выполнение через Handler
    createdQuest := casesteps.CreateQuestStep(...)
    
    // Проверка с помощью Handler assertions
    handlerAssertions.VerifyQuestFullMatch(...)
}
```

### Repository Tests (`repository_tests/`)
```go
// quest_repository_test.go
func (s *Suite) TestQuestRepository_Save_Success() {
    // Создание domain объекта
    quest := domain.NewQuest(...)
    
    // Сохранение через Repository
    savedQuest := s.TestDIContainer.QuestRepository.Save(...)
    
    // Проверка персистентности
    foundQuest := s.TestDIContainer.QuestRepository.GetByID(...)
}
```

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

### Выбор типа теста
1. **E2E** - для тестирования полных сценариев пользователя
2. **HTTP** - для тестирования API endpoints и валидации
3. **Handler** - для тестирования application логики без HTTP
4. **Repository** - для тестирования персистентности и БД

### Создание нового теста

1. **Выберите подходящую папку** в `tests/`
2. **Добавьте build tag** `//go:build integration`
3. **Используйте базовые компоненты** из `core/`

#### Пример E2E теста:
```go
//go:build integration

package quest_e2e_tests

import (
    "context"
    "testing"
    "github.com/stretchr/testify/suite"
    "quest-manager/tests/integration/tests"
)

func (s *Suite) TestNewE2EScenario() {
    ctx := context.Background()
    
    // Используйте case_steps
    quest := casesteps.CreateRandomQuestStep(...)
    
    // Используйте assertions
    e2eAssertions := assertions.NewQuestE2EAssertions(s.Assert())
    e2eAssertions.VerifyE2EFlow(...)
}
```

#### Пример HTTP теста:
```go
func (s *Suite) TestNewHTTPEndpoint() {
    // Подготовка данных через генераторы
    requestData := testdatagenerators.CustomQuestData(...)
    
    // HTTP запрос через case_steps
    response := casesteps.NewHTTPRequestStep(...)
    
    // Проверка через HTTP assertions
    httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
    httpAssertions.VerifyHTTPResponse(...)
}
```

### Расширение компонентов

#### Новые assertions:
Создайте файл в `core/assertions/` следуя паттерну:
```go
type NewFeatureAssertions struct {
    assert *assert.Assertions
}

func NewNewFeatureAssertions(assert *assert.Assertions) *NewFeatureAssertions {
    return &NewFeatureAssertions{assert: assert}
}
```

#### Новые case_steps:
Добавьте функции в существующие файлы или создайте новый в `core/case_steps/`

#### Новые генераторы данных:
Расширьте `quest_generator.go` или создайте новый генератор