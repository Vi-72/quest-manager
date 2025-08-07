# Tests Directory

Этот каталог содержит все тесты проекта Quest Manager, организованные по уровням и назначению.

## 📁 Структура тестов

```
tests/
├── domain/                    # 🏗️ Доменные (unit) тесты
│   ├── quest_test.go          # Тесты модели Quest
│   ├── location_test.go       # Тесты модели Location
│   ├── assign_quest_test.go   # Тесты бизнес-логики назначения
│   └── ...
├── contracts/                 # 🤝 Контрактные тесты
│   ├── mocks/                 # Mock реализации
│   └── *_contracts_test.go    # Тесты контрактов интерфейсов
├── integration/               # 🔗 Интеграционные тесты
│   ├── tests/                 # Группы интеграционных тестов
│   │   ├── quest_e2e_tests/   # E2E тесты квестов
│   │   ├── quest_http_tests/  # HTTP API тесты
│   │   ├── quest_handler_tests/ # Handler слой тесты
│   │   └── repository_tests/  # Repository тесты
│   ├── core/                  # Общие компоненты для интеграционных тестов
│   │   ├── assertions/        # Переиспользуемые проверки
│   │   ├── case_steps/        # Шаги тестирования
│   │   ├── storage/           # Прямой доступ к БД для тестов
│   │   └── test_data_generators/ # Генераторы тестовых данных
│   └── README.md             # Подробная документация интеграционных тестов
└── pkg/                      # 📦 Тесты для вспомогательных пакетов
    └── errs/                 # Тесты пакета ошибок
```

## 🧪 Типы тестов

### 1. **Domain Tests** (`tests/domain/`)
**Unit тесты доменной логики**
- ✅ Изолированные тесты бизнес-логики
- ✅ Быстрые и независимые
- ✅ Не требуют внешних зависимостей
- ✅ Тестируют правила домена, валидацию, агрегаты

```bash
# Запуск доменных тестов
make test-unit
# или
go test ./tests/domain -v
```

### 2. **Contract Tests** (`tests/contracts/`)
**Тесты контрактов интерфейсов**
- ✅ Проверяют соответствие реализаций интерфейсам
- ✅ Используют mock объекты
- ✅ Тестируют портыapplication слоя

```bash
# Запуск контрактных тестов  
go test ./tests/contracts -v
```

### 3. **Integration Tests** (`tests/integration/`)
**Тесты взаимодействия компонентов**
- ✅ Тесты с реальной базой данных PostgreSQL
- ✅ HTTP API тесты
- ✅ Repository слой тесты
- ✅ Handler и use case тесты
- ✅ E2E тесты полных сценариев

```bash
# Запуск всех интеграционных тестов
make test-integration
# или
go test -tags=integration ./tests/integration/... -v

# Запуск отдельных групп
go test -tags=integration ./tests/integration/tests/quest_http_tests -v
go test -tags=integration ./tests/integration/tests/repository_tests -v
```

### 4. **Package Tests** (`tests/pkg/`)
**Тесты вспомогательных пакетов**
- ✅ Тесты error handling
- ✅ Тесты утилит

## 🚀 Быстрый старт

### Запуск всех тестов
```bash
make test-all           # Все типы тестов
make test              # Unit + integration
```

### Запуск по типам
```bash
make test-unit         # Только unit тесты (быстро)
make test-integration  # Только integration тесты
make test-repository   # Только repository тесты
```

### Анализ покрытия
```bash
make coverage-check    # Быстрая проверка покрытия internal/ кода
make coverage-report   # Подробный отчет покрытия
make test-coverage     # Полное покрытие всех тестов
```

### Статистика тестов
```bash
make test-stats        # Подробная статистика по всем тестам
make test-stats-new    # Новая версия статистики
```

## 🛠️ Требования для запуска

### Для Unit тестов
- Go 1.21+
- Никаких дополнительных зависимостей

### Для Integration тестов
- PostgreSQL (через Docker Compose)
- Build tag: `-tags=integration`

```bash
# Запуск PostgreSQL
docker compose up -d postgres

# Создание тестовой БД (если нужно)
docker exec -it quest-manager-postgres-1 psql -U postgres -c "CREATE DATABASE quest_manager_test;"
```

## 📊 Покрытие кода

Цель: **>80%** покрытия для `internal/` кода

- ✅ **Domain layer**: >90% (бизнес-логика критична)
- ✅ **Application layer**: >85% (use cases и команды)
- ✅ **Adapters layer**: >70% (HTTP handlers, repositories)

Исключения из покрытия:
- `tests/` папки
- `cmd/` main функции
- Generated код

## 🔧 Добавление новых тестов

### Unit тесты
1. Создайте файл в `tests/domain/` с суффиксом `_test.go`
2. Импортируйте тестируемые пакеты из `internal/`
3. Используйте `testify/assert` для проверок

### Integration тесты
1. Создайте папку в `tests/integration/tests/`
2. Добавьте build tag `//go:build integration`
3. Используйте существующие компоненты из `core/`
4. Наследуйтесь от базовых test suites

Подробнее: [Integration Tests README](integration/README.md)

## 🎯 Best Practices

1. **Именование**: `TestFunction_Scenario_ExpectedResult`
2. **Структура**: Arrange → Act → Assert
3. **Изоляция**: Каждый тест независим
4. **Данные**: Используйте генераторы из `test_data_generators/`
5. **Скорость**: Unit тесты < 100ms, Integration < 1s
6. **Читаемость**: Тесты как живая документация

## 📈 CI/CD

Все тесты автоматически запускаются при:
- Pull requests
- Push в main ветку
- Nightly builds для статистики покрытия

Статус тестов отображается в README проекта.