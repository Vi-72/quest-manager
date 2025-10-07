# Quest Manager Service

HTTP-сервис для создания и управления квестами с поддержкой геолокаций.

## ✨ Основные возможности

- 🎯 **Управление квестами**: создание, назначение, изменение статуса
- 🔐 **JWT аутентификация**: защита всех API эндпоинтов через Bearer tokens
- 🗺️ **Геолокационный поиск**: поиск квестов по радиусу с точными расчетами
- 📍 **Гибридное хранение локаций**: денормализованные координаты + именованные локации
- ✅ **Продвинутая валидация**: многоуровневая система с детальными ошибками
- 🔄 **Domain Events**: отслеживание изменений в доменной модели
- 🏗️ **Clean Architecture**: четкое разделение слоев и ответственности
- ⚡ **Оптимизированная БД**: индексы для быстрого поиска
- 🚀 **style Container**: lazy initialization, context-aware dependencies
- 🔧 **Configuration-driven Middleware**: гибкая настройка через environment variables

## 🔧 Запуск

### 📦 Требования
- Go 1.23+
- PostgreSQL
- Quest Auth Service (gRPC) - для JWT аутентификации

### 🚀 Быстрый старт

1. **Настройка переменных окружения:**
```bash
cp config.example .env
# Отредактируйте .env файл под вашу конфигурацию
```

Обязательные переменные окружения:
```bash
HTTP_PORT=8080                          # Порт HTTP сервера
DB_HOST=localhost                       # PostgreSQL host
DB_PORT=5432                            # PostgreSQL port
DB_USER=postgres                        # Database user
DB_PASSWORD=secret                      # Database password
DB_NAME=quest_manager                   # Database name
DB_SSL_MODE=disable                     # SSL mode
EVENT_GOROUTINE_LIMIT=10               # Лимит горутин для событий
AUTH_GRPC=localhost:50051         # gRPC адрес Auth сервиса

# Middleware Configuration (опционально)
ENABLE_AUTH_MIDDLEWARE=true            # Включить аутентификацию
# Validation, Logging, Recovery - всегда включены
```

2. **Запуск:**
```bash
go run ./cmd/app
```

Сервер запускается на порту, указанном в переменной `HTTP_PORT` (по умолчанию 8080).

### 🔐 Аутентификация

Все API эндпоинты требуют JWT аутентификации. Добавьте токен в заголовок `Authorization`:

```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
     http://localhost:8080/api/v1/quests
```

**Коды ошибок аутентификации:**
- `401 Unauthorized` - невалидный, истекший или отсутствующий токен
- `403 Forbidden` - недостаточно прав (для будущих ролей)

### 🌐 API Endpoints

**Все эндпоинты требуют JWT аутентификации!**

- `GET /api/v1/quests` - Список всех квестов (с фильтрацией по статусу)
- `POST /api/v1/quests` - Создание нового квеста (возвращает location IDs)
- `GET /api/v1/quests/{quest_id}` - Получение квеста по ID (с валидацией UUID)
- `PATCH /api/v1/quests/{quest_id}/status` - Изменение статуса квеста
- `POST /api/v1/quests/{quest_id}/assign` - Назначение квеста пользователю
- `GET /api/v1/quests/assigned?user_id={id}` - Квесты назначенные пользователю
- `GET /api/v1/quests/search-radius` - Поиск квестов по радиусу

**Служебные эндпоинты (без аутентификации):**
- `GET /health` - Health check
- `GET /docs` - Swagger UI
- `GET /openapi.json` - OpenAPI спецификация

### 📖 Документация API

После запуска приложения доступна Swagger UI документация:
- Swagger UI: `http://localhost:8080/docs`
- OpenAPI JSON: `http://localhost:8080/openapi.json`

### 🏗️ Структура проекта

```
quest-manager/
├── cmd/                    # 🚀 Точка входа
│   ├── app/                # Главное приложение
│   ├── container.go        # DI контейнер
│   ├── build.go            # Build и валидация контейнера
│   ├── middlewares.go      # HTTP middleware
│   ├── router.go           # HTTP роутер
│   ├── closer.go           # Resource cleanup
│   └── config.go           # Конфигурация
├── internal/               # 🏗️ Основной код приложения
│   ├── adapters/           # Адаптеры (Hexagonal Architecture)
│   │   ├── in/http/        # HTTP handlers & middleware
│   │   │   ├── middleware/ # Auth & validation middleware
│   │   │   └── errors/     # Error handling (Problem Details)
│   │   └── out/            # Outbound adapters
│   │       ├── postgres/   # Репозитории БД
│   │       │   ├── questrepo/  # Quest repository
│   │       │   └── locationrepo/ # Location repository
│   │       └── client/auth/ # gRPC Auth client  
│   ├── core/               # Бизнес-логика (DDD)
│   │   ├── application/    # Use cases & handlers
│   │   ├── domain/         # Доменная модель
│   │   │   ├── quest/      # Quest aggregate
│   │   │   ├── location/   # Location aggregate
│   │   │   └── kernel/     # Shared value objects
│   │   └── ports/          # Интерфейсы
│   ├── generated/          # Сгенерированный код (OpenAPI)
│   └── pkg/                # Общие пакеты
│       ├── ddd/            # DDD building blocks
│       └── errs/           # Error types
├── tests/                  # 🧪 Все тесты проекта
│   ├── domain/             # Unit тесты доменной логики
│   ├── contracts/          # Контрактные тесты с моками
│   ├── integration/        # Интеграционные тесты
│   │   ├── tests/          # Группы тестов по слоям
│   │   └── core/           # Переиспользуемые компоненты
│   └── pkg/                # Тесты утилит
├── scripts/                # 📜 Скрипты для разработки и CI
│   ├── coverage-check.sh   # Быстрая проверка покрытия
│   ├── coverage-report.sh  # Детальный отчет покрытия
│   ├── test-stats.sh       # Статистика тестов
│   ├── test-stats-new.sh   # Новая статистика тестов
│   └── README.md           # Документация скриптов
├── .github/                # 🤖 GitHub Actions CI/CD
│   ├── workflows/          
│   │   └── ci.yml          # Основной CI pipeline
│   └── README.md           # Документация CI/CD
├── api/openapi/            # 📋 OpenAPI спецификация
├── configs/                # ⚙️ Конфигурационные файлы
├── .golangci.yml           # 🔍 Конфигурация линтера
├── .codecov.yml            # 📊 Конфигурация покрытия
└── Makefile                # 🛠️ Команды для разработки
```

### 🎯 Доменная модель

**Quest (Квест)** - Aggregate Root
- ID, Title, Description
- Difficulty (easy/medium/hard)
- Status (created/posted/assigned/in_progress/declined/completed)
- Target/Execution Location (координаты)
- Target/Execution Location IDs (ссылки на именованные локации)
- Equipment, Skills (списки)
- Creator, Assignee
- Timestamps
- Domain Events (QuestCreated, QuestAssigned, QuestStatusChanged)

**Location (Локация)** - Aggregate Root
- ID, Name (опционально), Address, Description
- Coordinate (GeoCoordinate)
- Timestamps
- Domain Events (LocationCreated, LocationUpdated)

**GeoCoordinate (Координаты)** - Value Object
- Latitude, Longitude
- Валидация диапазонов (-90..90, -180..180)
- Расчет расстояния (Haversine formula)
- Bounding box расчеты для оптимизации поиска

### 🗺️ Гибридное хранение локаций

Система использует **гибридный подход** для оптимального баланса производительности и гибкости:

1. **Денормализованные координаты** в таблице `quests`
   - Быстрый доступ для отображения и поиска
   - Всегда доступны даже без связанных локаций

2. **Именованные локации** в таблице `locations` 
   - Переиспользование популярных мест
   - Дополнительные метаданные (название, адрес, описание)
   - Опциональные FK в `quests.target_location_id` и `quests.execution_location_id`

```sql
-- Автоматически создаются при создании квеста
INSERT INTO locations (id, name, latitude, longitude, address, description)
VALUES (uuid, '', lat, lon, '', '');
```

## 🔍 Архитектура валидации

Система валидации построена на принципе **разделения ответственности** с **правильными HTTP кодами**:

### 📝 Уровни валидации

1. **Технические проверки** (OpenAPI middleware)
   - Форматы, обязательные поля, enum и диапазоны значений
   - Выполняются автоматически через `internal/adapters/in/http/middleware`
   - **Результат**: 400 Bad Request (Problem Details)

2. **Бизнес-правила** (доменная модель)
   - Enum значения (difficulty, status)
   - Бизнес-инварианты и переходы состояний
   - Доменная логика создания объектов
   - **Результат**: 400 Bad Request (DomainValidationError)

3. **Ресурсы** (application layer)
   - Существование объектов по ID
   - **Результат**: 404 Not Found (NotFoundError)

### 🚨 Обработка ошибок

```go
ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
    switch {
    case errors.As(err, &domainErr):
        problems.NewDomainValidationProblem(domainErr).WriteResponse(w)
    case errors.As(err, &notFoundErr):
        problems.NewNotFoundProblem(notFoundErr).WriteResponse(w)
    default:
        problems.NewBadRequest("Response error: " + err.Error()).WriteResponse(w)
    }
}
```

### 🔄 Процесс валидации

```go
// 1. HTTP Layer - технические проверки
validationmiddleware.Validate(r) // latitude/longitude format, ranges, required fields

// 2. Domain Layer - бизнес-правила  
quest, err := quest.NewQuest(dto.Title, dto.Difficulty, ...)
// difficulty enum, business invariants

// 3. Application Layer - ресурсы
quest, err := repository.GetByID(questID)
// existence checks
```

### ✅ Преимущества

- **Правильные HTTP коды**: 400 vs 404 vs 500
- **Четкое разделение**: техническая vs бизнес vs ресурсы
- **RFC 7807 Problem Details**: структурированные ошибки
- **Централизованная обработка**: middleware catch-all
- **Переиспользование**: общие функции валидации
- **Тестируемость**: независимые уровни

### 📋 Примеры ошибок

```json
// Техническая ошибка (400)
{
  "type": "bad-request",
  "title": "Bad Request", 
  "status": 400,
  "detail": "validation failed: field 'latitude' must be between -90 and 90"
}

// Бизнес-ошибка (400)
{
  "type": "bad-request", 
  "title": "Bad Request",
  "status": 400,
  "detail": "invalid status: must be one of 'created', 'posted', 'assigned'"
}

// Ресурс не найден (404)
{
  "type": "not-found",
  "title": "Not Found", 
  "status": 404,
  "detail": "quest with ID 'invalid-uuid' not found"
}
```

## ⚡ Производительность

### 🗂️ Индексы БД

```sql
-- Поиск по статусу
CREATE INDEX idx_quests_status ON quests(status);

-- Поиск по создателю/исполнителю  
CREATE INDEX idx_quests_creator ON quests(creator);
CREATE INDEX idx_quests_assignee ON quests(assignee);

-- Геопространственный поиск
CREATE INDEX idx_target_location ON quests(target_latitude, target_longitude);
CREATE INDEX idx_execution_location ON quests(execution_latitude, execution_longitude);

-- Локации
CREATE INDEX idx_locations_coords ON locations(latitude, longitude);
CREATE INDEX idx_locations_name ON locations(name);
```

### 🎯 Оптимизации поиска

1. **Bounding Box + Haversine**: сначала грубый поиск по прямоугольнику, затем точное расстояние
2. **Денормализация координат**: избегаем JOIN для частых запросов
3. **Композитные индексы**: оптимальны для multi-column поиска

## 🚀 Container Architecture

### 🏗️ Dependency Injection Container

Проект использует **Container** с современными паттернами:

#### **Lazy Initialization**
```go
// Зависимости создаются только при первом обращении
func (c *Container) GetAuthClient(ctx context.Context) ports.AuthClient {
    if c.authClient == nil {
        conn, err := grpc.NewClient(c.configs.AuthGRPC, ...)
        if err != nil {
            panic(fmt.Errorf("failed to create auth gRPC client: %w", err))
        }
        c.RegisterCloser(connCloser{conn})
        c.authClient = authclient.NewUserAuthClient(grpcClient)
    }
    return c.authClient
}
```

#### **Context-Aware Dependencies**
```go
// Все getter методы принимают context.Context
func (c *Container) GetAuthConn(ctx context.Context) *grpc.ClientConn
func (c *Container) GetAuthClient(ctx context.Context) ports.AuthClient
func (c *Container) GetQuestRepository(ctx context.Context) ports.QuestRepository
```

#### **Build Pattern**
```go
// Валидация и инициализация в отдельном методе
func (c *Container) Build(ctx context.Context) error {
    // Валидация конфигурации
    if c.configs.AuthGRPC != "" && c.configs.AuthClient != nil {
        return fmt.Errorf("both AuthGRPC and AuthClient cannot be set simultaneously")
    }
    
    // Eager validation для критических зависимостей
    if c.configs.AuthGRPC != "" {
        _ = c.GetAuthClient(ctx) // Trigger panic if fails
    }
    
    return nilCheck(c)
}
```

#### **Configuration-Driven Middleware**
```go
type MiddlewareConfig struct {
    EnableAuth       bool // Включает аутентификацию
    EnableValidation bool // Включает валидацию OpenAPI
    EnableLogging    bool // Включает логирование запросов
    EnableRecovery   bool // Включает recovery от паник
}

// Условная логика middleware
func (c *Container) Middlewares(swagger *openapi3.T) []func(http.Handler) http.Handler {
    if c.configs.Middleware.EnableAuth {
        if authClient := c.GetAuthClient(ctx); authClient != nil {
            authMW := httpmiddleware.NewAuthMiddleware(authClient)
            middlewares = append(middlewares, authMW.Auth)
            log.Printf("✅ Authentication middleware enabled")
        }
    }
    return middlewares
}
```

### 🎯 Преимущества

- ✅ **Lazy Loading**: зависимости создаются по требованию
- ✅ **Context Awareness**: все методы принимают context.Context
- ✅ **Panic on Critical Errors**: критические ошибки приводят к panic
- ✅ **Resource Management**: автоматическая регистрация closers
- ✅ **Configuration Flexibility**: middleware настраивается через env vars
- ✅ **Detailed Logging**: подробное логирование инициализации
- ✅ **Error Handling**: Build() возвращает error, getters panic

## 🔄 Domain-Driven Design

### 🏗️ Паттерны

- **Aggregate Root**: Quest, Location с инкапсуляцией бизнес-логики
- **Value Objects**: GeoCoordinate, BoundingBox  
- **Domain Events**: отслеживание изменений состояния
- **Unit of Work**: атомарные транзакции
- **Repository**: абстракция над хранилищем

### 📡 События

```go
// Автоматически создаются при изменениях
QuestCreated{ID, Title, CreatedAt, ...}
QuestAssigned{QuestID, UserID, AssignedAt, ...}  
QuestStatusChanged{QuestID, OldStatus, NewStatus, ...}

LocationCreated{ID, Coordinate, CreatedAt, ...}
LocationUpdated{ID, Coordinate, UpdatedAt, ...}
```

## 🚀 Генерация кода

Для регенерации HTTP сервера из OpenAPI:
```bash
make generate
# или
oapi-codegen -config configs/server.cfg.yaml api/openapi/openapi.yml
```

## 📚 Используемые библиотеки

- [Chi Router](https://github.com/go-chi/chi) - HTTP роутер
- [GORM](https://gorm.io/) - ORM для работы с БД
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - Генерация кода из OpenAPI
- [UUID](https://github.com/google/uuid) - Генерация UUID
- [gRPC](https://grpc.io/) - Интеграция с Auth сервисом
- [Quest Auth SDK](https://github.com/Vi-72/quest-auth) - gRPC SDK для аутентификации

## 🧪 Тестирование

### 📊 Покрытие кода: **75.6%**

![CI Status](https://github.com/Vi-72/quest-manager/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/Vi-72/quest-manager/branch/main/graph/badge.svg)](https://codecov.io/gh/Vi-72/quest-manager)

### 🎯 Результаты тестирования

#### **✅ Успешные тесты:**
- **Domain Tests**: 100% PASS - вся бизнес-логика работает корректно
- **Contract Tests**: 100% PASS - все интерфейсы и контракты соблюдены
- **Handler Tests**: 100% PASS - application layer работает стабильно

#### **⚠️ Проблемные тесты:**
- **HTTP Tests**: частично FAIL - проблемы с JSON unmarshaling и валидацией
- **E2E Tests**: 1 FAIL - создание квестов через API возвращает 400 вместо 201

**Примечание**: Проблемы в HTTP тестах не связаны с архитектурными изменениями Container - это существующие проблемы с HTTP слоем и валидацией.

### 🎯 Типы тестов

#### **Unit Tests** - Доменная логика
```bash
make test-unit          # Быстрые unit тесты
go test ./tests/domain -v
```

#### **Integration Tests** - Полный стек с PostgreSQL
```bash
make test-integration   # Все интеграционные тесты
go test -tags=integration ./tests/integration/... -v

# По группам:
go test -tags=integration ./tests/integration/tests/quest_e2e_tests -v      # E2E тесты
go test -tags=integration ./tests/integration/tests/quest_http_tests -v     # HTTP API
go test -tags=integration ./tests/integration/tests/quest_handler_tests -v  # Handlers
go test -tags=integration ./tests/integration/tests/repository_tests -v     # Repository
```

#### **Contract Tests** - Интерфейсы
```bash
go test ./tests/contracts -v
```

### 📈 Анализ покрытия

```bash
make coverage-check     # 🎯 Быстрая проверка покрытия internal/ кода
make coverage-report    # 📋 Подробный HTML отчет
make test-coverage      # 📊 Полное покрытие всех тестов
```

### 📊 Статистика тестов

```bash
make test-stats         # 📈 Подробная статистика по всем тестам  
make test-stats-new     # 📊 Новая версия статистики
```

### 🚀 Все тесты сразу

```bash
make test-all          # Unit + Integration + Contract тесты
make test              # Unit + Integration тесты
```

### 🎯 Качественные метрики

- ✅ **Domain Layer**: >90% покрытия (бизнес-критичная логика)
- ✅ **Application Layer**: >85% покрытия (use cases) 
- ✅ **Infrastructure Layer**: >70% покрытия (адаптеры)
- ✅ **Все тесты**: автоматически в CI/CD при каждом коммите

### 🔧 Требования для интеграционных тестов

```bash
# PostgreSQL через Docker
docker compose up -d postgres

# Создание тестовой БД (автоматически)
CREATE DATABASE quest_manager_test;
```

### 📁 Структура тестов

```
tests/
├── domain/                    # 🏗️ Unit тесты доменной логики
├── contracts/                 # 🤝 Тесты интерфейсов с моками  
├── integration/               # 🔗 Интеграционные тесты
│   ├── tests/                 # Группы тестов по слоям
│   │   ├── quest_e2e_tests/   # E2E полный цикл
│   │   ├── quest_http_tests/  # HTTP API тесты
│   │   ├── quest_handler_tests/ # Application handlers
│   │   └── repository_tests/  # Infrastructure repositories
│   └── core/                  # Переиспользуемые компоненты
│       ├── assertions/        # Специализированные проверки
│       ├── case_steps/        # Шаги тестирования
│       └── test_data_generators/ # Генераторы данных
└── pkg/                       # 📦 Тесты утилит
```

Подробнее: [Tests Documentation](tests/README.md)

## 🚀 CI/CD Pipeline

### 📋 GitHub Actions

Автоматическое тестирование при каждом push и pull request:

- ✅ **Unit Tests** - доменная логика и контрактные тесты
- ✅ **Integration Tests** - полный стек с PostgreSQL  
- ✅ **Coverage Report** - автоматическая отправка в Codecov
- ✅ **Linting** - проверка качества кода с golangci-lint
- ✅ **Build Check** - компиляция приложения

### 💬 Автоматические комментарии в PR

CI автоматически комментирует pull request'ы с:
- 📊 Актуальным покрытием кода
- 📈 Сравнением с предыдущими версиями  
- ✅ Статусом всех типов тестов
- 🎯 Соответствием целевым метрикам (>70%)

### 🎯 Качественные гейты

- **Минимальное покрытие**: 70% для `internal/` кода
- **Все тесты**: должны проходить без ошибок
- **Linting**: без нарушений стандартов кода
- **Build**: успешная компиляция

### 🔧 Скрипты тестирования

Все скрипты перенесены в [`scripts/`](scripts/) для удобства:

```bash
# Используются в CI и локально
make coverage-check     # scripts/coverage-check.sh
make test-stats         # scripts/test-stats.sh  
make test-stats-new     # scripts/test-stats-new.sh
make coverage-report    # scripts/coverage-report.sh
```

### 📈 Мониторинг качества

- **Codecov Dashboard**: история изменений покрытия
- **GitHub Actions**: статус всех сборок
- **PR Reviews**: блокировка merge при падении тестов

## 🔧 Разработка

Проект следует принципам **Clean Architecture** и **Domain-Driven Design**:

- **Domain Layer**: Богатая доменная модель с бизнес-правилами
- **Application Layer**: Use cases, CQRS handlers, Unit of Work
- **Infrastructure Layer**: Репозитории, внешние адаптеры
- **Ports & Adapters**: Инверсия зависимостей, тестируемость

### 🎯 Архитектурные решения

- **CQRS**: разделение команд и запросов
- **Упрощение**: удаление over-engineered структур
- **Event Sourcing Ready**: domain events для аудита
- **Hexagonal Architecture**: порты и адаптеры для изоляции
- **Database per Aggregate**: quest и location репозитории

## 🔄 Недавние изменения

### 🚀 Container Architecture Refactoring

**Дата**: Октябрь 2024

#### **Что изменилось:**
1. **Переименование**: `CompositionRoot` → `Container`
2. **Lazy Initialization**: зависимости создаются по требованию
3. **Context-Aware**: все getter методы принимают `context.Context`
4. **Build Pattern**: валидация вынесена в отдельный метод `Build()`
5. **Configuration-Driven Middleware**: гибкая настройка через env vars

#### **Удаленные файлы:**
- `cmd/auth_client_factory.go` - заменен на прямые getter методы
- `cmd/composition_root.go` - переименован в `container.go`

#### **Новые файлы:**
- `cmd/build.go` - валидация и инициализация контейнера
- `cmd/middlewares.go` - конфигурируемые HTTP middleware
- `cmd/router.go` - HTTP роутер с улучшенной логикой
- `cmd/closer.go` - управление ресурсами

#### **Новые возможности:**
```bash
# Middleware Configuration
ENABLE_AUTH_MIDDLEWARE=true
```

#### **Преимущества:**
- ✅ **Производительность**: lazy loading зависимостей
- ✅ **Гибкость**: настройка middleware через конфигурацию
- ✅ **Надежность**: panic на критических ошибках
- ✅ **Логирование**: подробная информация об инициализации
- ✅ **Тестируемость**: context-aware зависимости
