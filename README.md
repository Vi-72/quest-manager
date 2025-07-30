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