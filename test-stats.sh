#!/bin/bash

# Функция для запуска тестов с красивой статистикой
run_tests_with_stats() {
    local test_path="$1"
    local test_name="${2:-Tests}"
    
    echo "🧪 Running $test_name..."
    echo "═══════════════════════════════════════════════════════════════"
    
    # Запускаем тесты и сохраняем вывод
    local output=$(go test $test_path -v 2>&1)
    local exit_code=$?
    
    # Выводим полный вывод
    echo "$output"
    echo ""
    
    # Подсчитываем статистику
    local passed=$(echo "$output" | grep "\-\-\- PASS" | wc -l | tr -d ' ')
    local failed=$(echo "$output" | grep "\-\-\- FAIL" | wc -l | tr -d ' ')
    local skipped=$(echo "$output" | grep "\-\-\- SKIP" | wc -l | tr -d ' ')
    
    # Убеждаемся что переменные не пусты
    passed=${passed:-0}
    failed=${failed:-0}
    skipped=${skipped:-0}
    
    local total=$((passed + failed + skipped))
    
    # Определяем статус
    if [ $exit_code -eq 0 ]; then
        local status="✅ PASSED"
        local status_color="\033[32m" # Green
    else
        local status="❌ FAILED"
        local status_color="\033[31m" # Red
    fi
    
    # Выводим красивую статистику
    echo "📊 Test Results Summary:"
    echo "═══════════════════════════════════════════════════════════════"
    printf "${status_color}%s\033[0m\n" "$status"
    echo ""
    echo "📈 Statistics:"
    echo "   ✅ Passed:  $passed"
    echo "   ❌ Failed:  $failed"
    echo "   ⏭️  Skipped: $skipped"
    echo "   📝 Total:   $total"
    
    if [ $total -gt 0 ]; then
        local pass_rate=$((passed * 100 / total))
        echo "   📊 Success Rate: $pass_rate%"
    fi
    
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    
    return $exit_code
}

# Примеры использования:

echo "🚀 Test Statistics Script Ready!"
echo ""
echo "Usage examples:"
echo "  ./test-stats.sh domain          # Run domain tests"
echo "  ./test-stats.sh integration     # Run integration tests"
echo "  ./test-stats.sh http            # Run HTTP tests"
echo "  ./test-stats.sh handler         # Run handler tests"
echo "  ./test-stats.sh contracts       # Run contract tests"
echo "  ./test-stats.sh assign-quest    # Run assign quest tests"
echo "  ./test-stats.sh all             # Run all tests"
echo ""

# Обработка аргументов
case "${1:-help}" in
    "domain")
        run_tests_with_stats "./tests/domain/..." "Domain Tests"
        ;;
    "integration")
        run_tests_with_stats "./tests/integration/..." "Integration Tests"
        ;;
    "http")
        run_tests_with_stats "./tests/integration/cases/quest_http/..." "HTTP API Tests"
        ;;
    "handler")
        run_tests_with_stats "./tests/integration/cases/quest_handler/..." "Handler Tests"
        ;;
    "contracts")
        run_tests_with_stats "./tests/contracts/..." "Contract Tests"
        ;;
    "assign-quest")
        echo "🎯 Running Assign Quest Tests on All Layers..."
        echo ""
        
        run_tests_with_stats "./tests/domain/assign_quest_test.go ./tests/domain/quest_test.go" "Domain: Assign Quest"
        run_tests_with_stats "./tests/integration/cases/quest_handler/assign_quest_test.go ./tests/integration/cases/quest_handler/suite_test.go" "Handler: Assign Quest"
        run_tests_with_stats "./tests/integration/cases/quest_http/assign_quest_http_test.go ./tests/integration/cases/quest_http/suite_test.go" "HTTP: Assign Quest"
        ;;
    "all")
        run_tests_with_stats "./tests/..." "All Tests"
        ;;
    "help"|*)
        echo "Available commands listed above ☝️"
        ;;
esac