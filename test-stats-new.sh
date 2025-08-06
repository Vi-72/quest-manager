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
    
    # Возвращаем статистику через глобальные переменные для all команды
    export LAST_PASSED=$passed
    export LAST_FAILED=$failed
    export LAST_SKIPPED=$skipped
    export LAST_EXIT_CODE=$exit_code
    
    return $exit_code
}

# Функция для repository тестов с integration tag
run_repository_tests() {
    echo "🧪 Running Repository Tests with integration tag..."
    echo "═══════════════════════════════════════════════════════════════"
    
    # Запускаем тесты и сохраняем вывод
    output=$(go test -tags=integration ./tests/integration/tests/repository_tests -v 2>&1)
    exit_code=$?
    
    # Выводим полный вывод
    echo "$output"
    echo ""
    
    # Подсчитываем статистику
    passed=$(echo "$output" | grep "\-\-\- PASS" | wc -l | tr -d ' ')
    failed=$(echo "$output" | grep "\-\-\- FAIL" | wc -l | tr -d ' ')
    skipped=$(echo "$output" | grep "\-\-\- SKIP" | wc -l | tr -d ' ')
    
    # Убеждаемся что переменные не пусты
    passed=${passed:-0}
    failed=${failed:-0}
    skipped=${skipped:-0}
    
    total=$((passed + failed + skipped))
    
    # Определяем статус
    if [ $exit_code -eq 0 ]; then
        status="✅ PASSED"
        status_color="\033[32m"
    else
        status="❌ FAILED"
        status_color="\033[31m"
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
        pass_rate=$((passed * 100 / total))
        echo "   📊 Success Rate: $pass_rate%"
    fi
    
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    
    # Возвращаем статистику через глобальные переменные
    export LAST_PASSED=$passed
    export LAST_FAILED=$failed
    export LAST_SKIPPED=$skipped
    export LAST_EXIT_CODE=$exit_code
    
    return $exit_code
}

# Примеры использования:

echo "🚀 Test Statistics Script Ready!"
echo ""
echo "Usage examples:"
echo "  ./test-stats.sh domain          # Run domain tests"
echo "  ./test-stats.sh integration     # Run all integration tests"
echo "  ./test-stats.sh http            # Run HTTP API tests"
echo "  ./test-stats.sh handler         # Run handler tests"
echo "  ./test-stats.sh contracts       # Run contract tests"
echo "  ./test-stats.sh e2e              # Run E2E tests"
echo "  ./test-stats.sh repository      # Run repository tests"
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
        run_tests_with_stats "./tests/integration/tests/quest_http_tests/..." "HTTP API Tests"
        ;;
    "handler")
        run_tests_with_stats "./tests/integration/tests/quest_handler_tests/..." "Handler Tests"
        ;;
    "contracts")
        run_tests_with_stats "./tests/contracts/..." "Contract Tests"
        ;;
    "e2e")
        run_tests_with_stats "./tests/integration/tests/quest_e2e_tests/..." "E2E Tests"
        ;;
    "repository")
        run_repository_tests
        ;;
    "all")
        echo "🎯 Running All Test Suites..."
        echo ""
        
        # Переменные для общего подсчета
        total_passed=0
        total_failed=0
        total_skipped=0
        overall_exit_code=0
        
        # Запускаем каждый набор тестов и собираем статистику
        echo ""
        run_tests_with_stats "./tests/domain/..." "Domain Tests"
        total_passed=$((total_passed + LAST_PASSED))
        total_failed=$((total_failed + LAST_FAILED))
        total_skipped=$((total_skipped + LAST_SKIPPED))
        if [ $LAST_EXIT_CODE -ne 0 ]; then
            overall_exit_code=1
        fi
        
        echo ""
        run_tests_with_stats "./tests/contracts/..." "Contract Tests"
        total_passed=$((total_passed + LAST_PASSED))
        total_failed=$((total_failed + LAST_FAILED))
        total_skipped=$((total_skipped + LAST_SKIPPED))
        if [ $LAST_EXIT_CODE -ne 0 ]; then
            overall_exit_code=1
        fi
        
        echo ""
        run_tests_with_stats "./tests/integration/tests/quest_http_tests/..." "HTTP API Tests"
        total_passed=$((total_passed + LAST_PASSED))
        total_failed=$((total_failed + LAST_FAILED))
        total_skipped=$((total_skipped + LAST_SKIPPED))
        if [ $LAST_EXIT_CODE -ne 0 ]; then
            overall_exit_code=1
        fi
        
        echo ""
        run_tests_with_stats "./tests/integration/tests/quest_handler_tests/..." "Handler Tests"
        total_passed=$((total_passed + LAST_PASSED))
        total_failed=$((total_failed + LAST_FAILED))
        total_skipped=$((total_skipped + LAST_SKIPPED))
        if [ $LAST_EXIT_CODE -ne 0 ]; then
            overall_exit_code=1
        fi
        
        echo ""
        run_tests_with_stats "./tests/integration/tests/quest_e2e_tests/..." "E2E Tests"
        total_passed=$((total_passed + LAST_PASSED))
        total_failed=$((total_failed + LAST_FAILED))
        total_skipped=$((total_skipped + LAST_SKIPPED))
        if [ $LAST_EXIT_CODE -ne 0 ]; then
            overall_exit_code=1
        fi
        
        echo ""
        run_repository_tests
        total_passed=$((total_passed + LAST_PASSED))
        total_failed=$((total_failed + LAST_FAILED))
        total_skipped=$((total_skipped + LAST_SKIPPED))
        if [ $LAST_EXIT_CODE -ne 0 ]; then
            overall_exit_code=1
        fi
        
        # ОБЩАЯ ИТОГОВАЯ СТАТИСТИКА
        overall_total=$((total_passed + total_failed + total_skipped))
        if [ $overall_exit_code -eq 0 ]; then
            overall_status="✅ ALL TESTS PASSED"
            overall_color="\033[32m" # Green
        else
            overall_status="❌ SOME TESTS FAILED"
            overall_color="\033[31m" # Red
        fi
        
        echo ""
        echo "🏆 OVERALL TEST RESULTS"
        echo "═══════════════════════════════════════════════════════════════"
        printf "${overall_color}%s\033[0m\n" "$overall_status"
        echo ""
        echo "📈 Overall Statistics:"
        echo "   ✅ Total Passed:  $total_passed"
        echo "   ❌ Total Failed:  $total_failed"
        echo "   ⏭️  Total Skipped: $total_skipped"
        echo "   📝 Grand Total:   $overall_total"
        if [ $overall_total -gt 0 ]; then
            overall_pass_rate=$((total_passed * 100 / overall_total))
            echo "   📊 Overall Success Rate: $overall_pass_rate%"
        fi
        echo "═══════════════════════════════════════════════════════════════"
        
        exit $overall_exit_code
        ;;
    "help"|*)
        echo "Available commands listed above ☝️"
        ;;
esac