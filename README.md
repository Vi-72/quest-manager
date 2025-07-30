# Quest Manager Service

HTTP-сервис для создания и управления квестами.


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

- `GET /api/v1/quests` - Список всех квестов
- `POST /api/v1/quests` - Создание нового квеста
- `GET /api/v1/quests/{quest_id}` - Получение квеста по ID
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
│   │   ├── in/http/        # HTTP handlers
│   │   └── out/postgres/   # Репозитории БД
│   ├── core/               # Бизнес-логика
│   │   ├── application/    # Use cases
│   │   ├── domain/         # Доменная модель
│   │   └── ports/          # Интерфейсы
│   ├── generated/          # Сгенерированный код
│   └── pkg/                # Общие пакеты
│       └── validations/    # Система валидации
├── api/openapi/            # OpenAPI спецификация
└── configs/                # Конфигурационные файлы
```

### 🎯 Доменная модель

**Quest (Квест)**
- ID, Title, Description
- Difficulty (easy/medium/hard)
- Status (created/posted/assigned/in_progress/declined/completed)
- Target/Execution Location (координаты)
- Equipment, Skills (списки)
- Creator, Assignee
- Timestamps

**GeoCoordinate (Координаты)**
- Latitude, Longitude
- Валидация диапазонов
- Расчет расстояния (Haversine formula)

## 🔍 Архитектура валидации

Система валидации построена на принципе **разделения ответственности**:

### 📝 Уровни валидации

1. **Технические проверки** (`internal/pkg/validations/`)
   - Формат данных (не пустые строки, диапазоны координат)
   - Синтаксис (UUID формат, числовые значения)
   - Безопасность (размеры полей, специальные символы)

2. **Бизнес-правила** (доменная модель)
   - Соответствие enum значениям (difficulty: easy/medium/hard)
   - Бизнес-логика создания объектов
   - Инварианты доменной модели

### 📁 Структура валидации

```
internal/pkg/validations/
├── common.go           # Базовые типы и общие функции
├── coordinates.go      # Валидация и конвертация координат
├── create_quest.go     # Валидация создания квеста  
└── assign_quest.go     # Валидация назначения квеста
```

### 🔄 Процесс валидации

```go
// 1. HTTP Layer - технические проверки
validatedData, err := validations.ValidateCreateQuestRequest(request.Body)

// 2. Domain Layer - бизнес-правила
quest, err := quest.NewQuest(
    validatedData.Title,
    validatedData.Description, 
    validatedData.Difficulty, // строка валидируется в домене
    // ...
)
```

### ✅ Преимущества

- **Четкое разделение**: техническая vs бизнес валидация
- **Переиспользование**: общие функции в `common.go`
- **Тестируемость**: независимые уровни валидации
- **Детальные ошибки**: RFC 7807 Problem Details
- **Типобезопасность**: строгая типизация на каждом уровне

### 📋 Примеры валидации

```json
// Техническая ошибка
{
  "type": "bad-request",
  "title": "Bad Request", 
  "status": 400,
  "detail": "validation failed: field 'title' is required and cannot be empty"
}

// Бизнес-ошибка
{
  "type": "bad-request",
  "title": "Bad Request",
  "status": 400, 
  "detail": "invalid difficulty: must be one of 'easy', 'medium', 'hard'"
}
```

## 🚀 Генерация кода

Для регенерации HTTP сервера из OpenAPI:
```bash
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

Проект следует принципам Clean Architecture:
- **Domain Layer**: Доменная модель и бизнес-правила
- **Application Layer**: Use cases и сценарии
- **Infrastructure Layer**: Адаптеры для внешних систем
- **Ports & Adapters**: Инверсия зависимостей