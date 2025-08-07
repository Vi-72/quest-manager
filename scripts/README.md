# Scripts Directory

Этот каталог содержит все вспомогательные скрипты проекта Quest Manager.

## Доступные скрипты

### 📊 Тестирование и статистика

#### `test-stats.sh`
Запускает тесты с подробной статистикой и красивым выводом.
```bash
# Прямой запуск
./scripts/test-stats.sh

# Через Makefile
make test-stats
```

#### `test-stats-new.sh`
Новая версия скрипта статистики тестов.
```bash
# Прямой запуск
./scripts/test-stats-new.sh

# Через Makefile
make test-stats-new
```

### 📈 Покрытие кода

#### `coverage-check.sh`
Проверяет покрытие только для internal/ кода, исключая tests/.
Показывает итоговый процент покрытия одним числом.
```bash
# Прямой запуск
./scripts/coverage-check.sh

# Через Makefile
make coverage-check
```

#### `coverage-report.sh`
Генерирует подробный отчет о покрытии кода.
```bash
# Прямой запуск
./scripts/coverage-report.sh

# Через Makefile
make coverage-report
```

## Запуск через Makefile

Все скрипты доступны через удобные цели в Makefile:

```bash
make test-stats         # Статистика тестов
make test-stats-new     # Новая статистика тестов
make coverage-check     # Быстрая проверка покрытия
make coverage-report    # Подробный отчет покрытия
```

## Разрешения

Makefile автоматически устанавливает права на выполнение для скриптов перед запуском.
Если запускаете скрипты напрямую, убедитесь что они исполняемые:

```bash
chmod +x scripts/*.sh
```