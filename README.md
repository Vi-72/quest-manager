# Quest Manager Service

HTTP-сервис для создания и управления квестами с поддержкой геолокаций.

## ✨ Основные возможности

- 🎯 **Управление квестами**: создание, назначение, изменение статуса
- 🗺️ **Геолокационный поиск**: поиск квестов по радиусу с точными расчетами
- 📍 **Гибридное хранение локаций**: денормализованные координаты + именованные локации
- ✅ **Продвинутая валидация**: многоуровневая система с детальными ошибками
- 🔄 **Domain Events**: отслеживание изменений в доменной модели
- 🏗️ **Clean Architecture**: четкое разделение слоев и ответственности
- ⚡ **Оптимизированная БД**: индексы для быстрого поиска

## 🔧 Запуск

### 📦 Требования
- Go 1.23+
- PostgreSQL

### 🚀 Быстрый старт

1. **Настройка переменных окружения:**
```bash
cp config.example .env
# Отредактируйте .env файл под вашу конфигурацию
```

2. **Запуск:**
```bash
go run ./cmd/app
```

Сервер запускается на порту, указанном в переменной `HTTP_PORT` (по умолчанию 8080).

### 🌐 API Endpoints

- `GET /api/v1/quests` - Список всех квестов (с фильтрацией по статусу)
- `POST /api/v1/quests` - Создание нового квеста (возвращает location IDs)
- `GET /api/v1/quests/{quest_id}` - Получение квеста по ID (с валидацией UUID)
- `PATCH /api/v1/quests/{quest_id}/status` - Изменение статуса квеста
- `POST /api/v1/quests/{quest_id}/assign` - Назначение квеста пользователю
- `GET /api/v1/quests/assigned?user_id={id}` - Квесты назначенные пользователю
- `GET /api/v1/quests/search-radius` - Поиск квестов по радиусу

### 📖 Документация API

После запуска приложения доступна Swagger UI документация:
- Swagger UI: `http://localhost:8080/docs`
- OpenAPI JSON: `http://localhost:8080/openapi.json`

### 🏗️ Структура проекта

```
quest-manager/
├── cmd/                    # Точка входа
│   ├── app/                # Главное приложение
│   ├── composition_root.go # DI контейнер
│   └── config.go           # Конфигурация
├── internal/
│   ├── adapters/           # Адаптеры
│   │   ├── in/http/        # HTTP handlers & validations
│   │   └── out/postgres/   # Репозитории БД
│   │       ├── questrepo/  # Quest repository
│   │       └── locationrepo/ # Location repository  
│   ├── core/               # Бизнес-логика
│   │   ├── application/    # Use cases & handlers
│   │   ├── domain/         # Доменная модель
│   │   │   ├── quest/      # Quest aggregate
│   │   │   ├── location/   # Location aggregate
│   │   │   └── kernel/     # Shared value objects
│   │   └── ports/          # Интерфейсы
│   ├── generated/          # Сгенерированный код
│   └── pkg/                # Общие пакеты
│       ├── ddd/            # DDD building blocks
│       └── errs/           # Error types
├── api/openapi/            # OpenAPI спецификация
└── configs/                # Конфигурационные файлы
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

1. **Технические проверки** (`internal/adapters/in/http/validations/`)
   - Формат данных (UUID, координаты, не пустые строки)
   - Синтаксис и диапазоны значений
   - Безопасность (размеры полей)
   - **Результат**: 400 Bad Request

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
// Кастомные типы ошибок
type DomainValidationError struct { Field, Message string }
type NotFoundError struct { Resource, ID string }

// Централизованная обработка в middleware
ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
    switch err := err.(type) {
    case *validations.ValidationError:
        // 400 Bad Request
    case *errs.DomainValidationError:
        // 400 Bad Request  
    case *errs.NotFoundError:
        // 404 Not Found
    default:
        // 500 Internal Server Error
    }
}
```

### 📁 Структура валидации

```
internal/adapters/in/http/validations/
├── common.go           # Базовые типы и общие функции
├── coordinates.go      # Валидация и конвертация координат
├── create_quest.go     # Валидация создания квеста  
├── assign_quest.go     # Валидация назначения квеста
├── change_quest_status.go # Валидация смены статуса
└── error_converters.go # Конвертация ошибок в Problem Details
```

### 🔄 Процесс валидации

```go
// 1. HTTP Layer - технические проверки
validatedData, err := validations.ValidateCreateQuestRequest(request.Body)
// latitude/longitude format, ranges, required fields

// 2. Domain Layer - бизнес-правила  
quest, err := quest.NewQuest(validatedData.Title, validatedData.Difficulty, ...)
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

## 🧪 Тестирование

```bash
go test ./...
```

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