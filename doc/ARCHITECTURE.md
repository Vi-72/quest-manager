# Quest Manager - Architecture Documentation

## 🏗️ Обзор архитектуры

Quest Manager построен на принципах **Clean Architecture** с использованием **Container Architecture** по образцу nfactors. Проект следует принципам **Domain-Driven Design (DDD)** и **Dependency Injection (DI)**.

## 📁 Структура проекта

```
quest-manager/
├── cmd/                           # 🚀 Точка входа и DI Container
│   ├── app/main.go               # Главное приложение
│   ├── container.go              # DI Container (nfactors-style)
│   ├── build.go                  # Build и валидация контейнера
│   ├── middlewares.go            # HTTP middleware
│   ├── router.go                 # HTTP роутер
│   ├── closer.go                 # Resource cleanup
│   └── config.go                 # Конфигурация
├── internal/
│   ├── adapters/
│   │   ├── in/http/              # HTTP адаптеры (входящие)
│   │   └── out/                  # Внешние адаптеры (исходящие)
│   │       ├── client/auth/      # Auth gRPC клиент + фабрика
│   │       └── postgres/         # PostgreSQL репозитории
│   ├── core/
│   │   ├── application/          # Application layer
│   │   │   ├── usecases/
│   │   │   │   ├── commands/     # Command handlers (CQRS)
│   │   │   │   └── queries/      # Query handlers (CQRS)
│   │   ├── domain/               # Domain layer
│   │   │   └── model/            # Domain models
│   │   └── ports/                # Ports (интерфейсы)
│   └── pkg/                      # Общие пакеты
└── tests/                        # Тесты всех уровней
```

## 🎯 Архитектурные принципы

### 1. Clean Architecture
- **Domain Layer** - бизнес-логика, не зависит от внешних слоев
- **Application Layer** - use cases, координирует domain и infrastructure
- **Infrastructure Layer** - внешние зависимости (БД, HTTP, gRPC)
- **Interface Layer** - HTTP handlers, API endpoints

### 2. Container Architecture (nfactors-style)
- **Lazy Initialization** - зависимости создаются по требованию
- **Context-Aware** - все getter методы принимают `context.Context`
- **Factory Pattern** - создание клиентов через фабрики
- **Resource Management** - автоматическая очистка ресурсов

### 3. Domain-Driven Design
- **Aggregate Root** - Quest, Location с инкапсуляцией бизнес-логики
- **Value Objects** - GeoCoordinate, BoundingBox
- **Domain Events** - отслеживание изменений состояния
- **Repository Pattern** - абстракция над хранилищем

## 🔧 Компоненты системы

### Container (DI Compose)

```go
type Container struct {
    configs Config
    db      *gorm.DB
    
    // Lazy initialized dependencies
    unitOfWork     ports.UnitOfWork
    eventPublisher ports.EventPublisher
    authClient     ports.AuthClient
    
    // Resource cleanup
    closers []Closer
}
```

**Принципы:**
- **Lazy Initialization** - зависимости создаются при первом обращении
- **Context-Aware** - все методы принимают `context.Context`
- **Resource Management** - автоматическая регистрация closers
- **Factory Delegation** - создание клиентов делегируется фабрикам

### Auth Factory

```go
type Factory struct {
    Addr   string
    Client ports.AuthClient
}

func (f *Factory) Create(ctx context.Context) (ports.AuthClient, *grpc.ClientConn, error)
```

**Возможности:**
- **Mock Support** - поддержка mock клиентов для тестов
- **Real gRPC** - создание реальных gRPC соединений
- **Modern API** - использует `grpc.NewClient()` вместо устаревшего `Dial`
- **Error Handling** - корректная обработка ошибок

### Middleware Configuration

```go
type MiddlewareConfig struct {
    EnableAuth bool  // Только auth настраивается
    // Validation, Logging, Recovery - всегда включены
}
```

**Middleware:**
- **Authentication** - JWT токены через gRPC
- **Validation** - OpenAPI валидация (всегда включена)
- **Logging** - HTTP логирование (всегда включено)
- **Recovery** - обработка паник (всегда включено)

## 🔄 Потоки данных

### 1. HTTP Request Flow

```
HTTP Request → Router → Middleware → Handler → Use Case → Domain → Repository → Database
                ↓
HTTP Response ← Mapper ← Use Case ← Domain ← Repository ← Database
```

### 2. Command Flow (CQRS)

```
HTTP POST → CreateQuestHandler → CreateQuestCommand → Quest Aggregate → Repository → Events
```

### 3. Query Flow (CQRS)

```
HTTP GET → ListQuestsHandler → ListQuestsQuery → Repository → Quest[] → HTTP Response
```

## 🧪 Тестирование

### Типы тестов

1. **Domain Tests** - бизнес-логика без внешних зависимостей
2. **Contract Tests** - интерфейсы между слоями
3. **Integration Tests** - полный цикл с базой данных
4. **HTTP Tests** - API endpoints
5. **E2E Tests** - end-to-end сценарии

### Изоляция тестов

- **Database Cleanup** - очистка БД между тестами
- **Mock Clients** - использование mock auth клиентов
- **Test Containers** - изолированные DI контейнеры
- **Event Processing** - ожидание завершения async операций

## 🚀 Запуск и конфигурация

### Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=quest_manager
DB_SSL_MODE=disable

# Application
HTTP_PORT=8080
EVENT_GOROUTINE_LIMIT=10

# Auth Service
AUTH_GRPC=localhost:50051

# Middleware (optional)
ENABLE_AUTH_MIDDLEWARE=true
```

### Команды

```bash
# Запуск приложения
go run ./cmd/app

# Тесты
go test ./tests/... -p 1

# Линтер
golangci-lint run

# Сборка
go build ./cmd/...
```

## 📊 Метрики и мониторинг

### Покрытие тестами
- **Domain Tests**: 60 тестов ✅
- **Contract Tests**: 61 тест ✅
- **Integration Tests**: 30 тестов ✅
- **HTTP Tests**: 57 тестов ✅
- **E2E Tests**: 4 теста ✅

**Общее покрытие: ~200+ тестов (100% PASS)**

### Производительность
- **Lazy Initialization** - быстрый старт приложения
- **Connection Pooling** - эффективное использование БД соединений
- **Event Processing** - асинхронная обработка событий
- **Resource Management** - автоматическая очистка ресурсов

## 🔮 Будущие улучшения

### Планируемые изменения
1. **Metrics & Monitoring** - Prometheus метрики
2. **Distributed Tracing** - OpenTelemetry
3. **Rate Limiting** - защита от DDoS
4. **Caching** - Redis для кеширования
5. **Message Queue** - асинхронная обработка событий

### Архитектурные улучшения
1. **Event Sourcing** - полная история изменений
2. **CQRS with Read Models** - оптимизированные read модели
3. **Microservices** - разделение на отдельные сервисы
4. **API Gateway** - централизованная точка входа

## 📚 Дополнительные ресурсы

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)
- [Dependency Injection](https://martinfowler.com/articles/injection.html)

---

**Версия документации**: 1.0  
**Дата обновления**: Октябрь 2024  
**Статус**: Production Ready ✅
