#!/bin/bash

# Скрипт для проверки покрытия только internal/ папки
# Исключает tests/ из подсчета покрытия
# Показывает итоговый процент покрытия одним числом

echo "🎯 Checking Coverage for Internal Code Only..."
echo "═══════════════════════════════════════════════════════════════"

# Запускаем все тесты с покрытием только для internal/
# Исключаем tests/ папку из подсчета покрытия
# Включаем repository тесты с build tag integration
go test -p 1 -count=1 -tags=integration -coverprofile=internal_coverage.out -coverpkg=./internal/... ./tests/... 2>/dev/null

if [ -f "internal_coverage.out" ]; then
    # Получаем итоговый процент покрытия
    total_result=$(go tool cover -func=internal_coverage.out | tail -1)
    total_percent=$(echo "$total_result" | awk '{print $NF}' | sed 's/%//')
    
    # Определяем цвет
    if (( $(echo "$total_percent >= 80" | bc -l) )); then
        color="\033[32m" # Green for 80%+
    elif (( $(echo "$total_percent >= 60" | bc -l) )); then
        color="\033[33m" # Yellow for 60-79%
    else
        color="\033[31m" # Red for <60%
    fi
    
    echo ""
    echo -e "📊 Internal Code Coverage: ${color}${total_percent}%\033[0m"
    echo ""
    echo "✅ Coverage calculation includes:"
    echo "   • internal/domain/"
    echo "   • internal/application/"
    echo "   • internal/adapters/"
    echo ""
    echo "❌ Coverage calculation excludes:"
    echo "   • tests/"
    echo "   • cmd/"
    echo "   • configs/"
    echo ""
    
    # Показываем итоговое число для скриптов
    echo "FINAL_COVERAGE_PERCENT: ${total_percent}"
    
    # Cleanup
    rm -f internal_coverage.out
    
    exit 0
else
    echo "❌ Failed to generate coverage report"
    exit 1
fi