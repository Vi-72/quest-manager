# GitHub Actions CI/CD

Этот каталог содержит GitHub Actions workflows для автоматизации тестирования, сборки и релизов Quest Manager.

## 🚀 Workflows

### `ci.yml` - Основной CI Pipeline
**Запускается:** При каждом push и pull request в main

**Задачи:**
- 🧪 **Unit Tests** - Доменные и контрактные тесты
- 🔗 **Integration Tests** - Тесты с PostgreSQL
- 📊 **Coverage Report** - Генерация и отправка покрытия в Codecov
- 🔍 **Linting** - Проверка кода с golangci-lint
- 🔨 **Build** - Компиляция приложения

**Особенности:**
- Автоматические комментарии с покрытием в PR
- Проверка минимального порога покрытия (70%)
- Использование наших скриптов из `scripts/`
- Кэширование Go modules для ускорения

### `nightly.yml` - Ночные тесты
**Запускается:** Каждую ночь в 2:00 UTC + вручную

**Задачи:**
- 🌙 **Расширенное тестирование** - Полный набор тестов
- 🚀 **Benchmark тесты** - Проверка производительности
- 📋 **Детальные отчеты** - coverage-report.sh
- 🔒 **Security scan** - Gosec и vulnerability check
- 🚨 **Уведомления** - Создание issue при падении

### `release.yml` - Релизы
**Запускается:** При создании git tag `v*`

**Задачи:**
- ✅ **Полное тестирование** - Все типы тестов
- 🔨 **Мультиплатформенная сборка** - Linux, macOS, Windows
- 📦 **GitHub Release** - Автоматическое создание релиза
- 📋 **Release Notes** - Генерация с метриками покрытия

## 📊 Интеграции

### Codecov
- Автоматическая отправка покрытия
- Комментарии в PR с изменениями покрытия
- Конфигурация в `.codecov.yml`

### golangci-lint
- Статический анализ кода
- Настройки в `.golangci.yml`
- 20+ включенных линтеров

## 🎯 Качественные гейты

### CI Pipeline
- ✅ Unit tests должны проходить
- ✅ Integration tests должны проходить
- ✅ Покрытие ≥70% для внутреннего кода
- ✅ Linting без ошибок
- ✅ Успешная сборка

### Nightly
- 📊 Детальные метрики покрытия
- 🚀 Benchmark тесты
- 🔒 Security сканирование
- 📈 Исторические данные

### Release
- 🎯 Покрытие ≥80% рекомендуется
- ✅ Все тесты проходят
- 🔨 Мультиплатформенные бинарники
- 📋 Автоматические release notes

## 🔧 Локальная эмуляция

### Запуск как в CI
```bash
# Unit тесты как в CI
go test ./tests/domain -v -race -coverprofile=unit_coverage.out

# Integration тесты как в CI
go test -tags=integration -v -race -p 1 -count=1 ./tests/integration/...

# Linting как в CI
golangci-lint run --timeout=5m

# Сборка как в CI
make build
```

### Проверка покрытия
```bash
# Быстрая проверка (как в CI)
make coverage-check

# Детальный отчет (как в nightly)
make coverage-report
```

## 📈 Мониторинг

### PR Checks
- Статус всех проверок отображается в PR
- Автоматические комментарии с покрытием
- Блокировка merge при падении тестов

### Metrics Dashboard
- Codecov dashboard для трендов покрытия
- GitHub Actions для истории билдов
- Issues для отслеживания nightly failures

## 🚨 Troubleshooting

### Частые проблемы
1. **PostgreSQL connection** - Проверьте health checks
2. **Coverage threshold** - Убедитесь что покрытие ≥70%
3. **Lint errors** - Запустите `golangci-lint run` локально
4. **Flaky tests** - Проверьте nightly reports

### Debug режим
Добавьте в workflow для отладки:
```yaml
- name: Debug info
  run: |
    echo "Go version: $(go version)"
    echo "PostgreSQL status: $(pg_isready -h localhost -p 5432 -U postgres)"
    echo "Environment: $DATABASE_URL"
```

## 🔄 Обновления

При изменении workflow'ов:
1. Тестируйте локально где возможно
2. Создавайте PR для review
3. Проверяйте на test branch сначала
4. Мониторьте первые запуски после изменений

---

*Все workflows используют наши существующие скрипты из `scripts/` для консистентности между локальной разработкой и CI.*