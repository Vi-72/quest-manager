#!/bin/bash

# Функция для анализа покрытия с красивой статистикой
run_coverage_analysis() {
    local test_path="$1"
    local test_name="${2:-Coverage Analysis}"
    local coverage_file="coverage_${3:-temp}.out"
    
    echo "📊 Running $test_name..."
    echo "═══════════════════════════════════════════════════════════════"
    
    # Запускаем тесты с покрытием
    local output=$(go test -p 1 -count=1 -coverprofile="$coverage_file" -coverpkg=./internal/... $test_path 2>&1)
    local exit_code=$?
    
    if [ $exit_code -ne 0 ]; then
        echo "❌ FAILED to run tests for coverage analysis"
        echo "$output"
        return $exit_code
    fi
    
    # Извлекаем общий процент покрытия
    local total_coverage=$(echo "$output" | grep -o "coverage: [0-9.]*%" | tail -1 | grep -o "[0-9.]*")
    
    if [ -z "$total_coverage" ]; then
        echo "⚠️  No coverage data available for $test_name"
        return 0
    fi
    
    # Определяем цвет для общего покрытия
    local coverage_color="\033[31m" # Default to Red
    if (( $(echo "$total_coverage >= 80" | bc -l) )); then
        coverage_color="\033[32m" # Green for 80%+
    elif (( $(echo "$total_coverage >= 60" | bc -l) )); then
        coverage_color="\033[33m" # Yellow for 60-79%
    fi
    
    echo ""
    echo -e "📈 Overall Coverage: ${coverage_color}${total_coverage}%\033[0m"
    echo ""
    
    # Генерируем детальный отчет по функциям
    if [ -f "$coverage_file" ]; then
        echo "🔍 Detailed Function Coverage:"
        echo "─────────────────────────────────────────────────────────────"
        
        # Парсим покрытие по модулям
        local func_coverage=$(go tool cover -func="$coverage_file" 2>/dev/null)
        
        if [ -n "$func_coverage" ]; then
            # Группируем по пакетам и показываем ключевые метрики (только internal/)
            echo "$func_coverage" | grep "/internal/" | \
            awk '{
                # Извлекаем процент покрытия
                coverage = $NF
                gsub(/%/, "", coverage)
                
                # Определяем пакет
                if ($0 ~ /\/http\//) package = "🌐 HTTP"
                else if ($0 ~ /\/domain\//) package = "🏗️ Domain"
                else if ($0 ~ /\/application\//) package = "🔧 Application"
                else if ($0 ~ /\/postgres\//) package = "💾 Repository"
                else package = "📦 Other"
                
                # Определяем цвет (принудительно преобразуем в число)
                coverage_num = coverage + 0
                if (coverage_num >= 80) color = "\033[32m"      # Green
                else if (coverage_num >= 60) color = "\033[33m" # Yellow
                else color = "\033[31m"                     # Red
                
                # Извлекаем имя функции
                func_name = $2
                gsub(/.*\//, "", func_name)
                
                printf "%-15s %s%-40s %6.1f%%\033[0m\n", package, color, func_name, coverage
            }' | sort -k1,1 -k3,3nr | head -20
            
            echo ""
            echo "📊 Coverage by Layer:"
            echo "─────────────────────────────────────────────────────────────"
            
            # Группируем статистику по слоям (только internal/)
            echo "$func_coverage" | grep "/internal/" | awk '
            BEGIN {
                layers["http"] = "🌐 HTTP Layer      "
                layers["domain"] = "🏗️ Domain Layer    "
                layers["application"] = "🔧 Application Layer"
                layers["postgres"] = "💾 Repository Layer"
            }
            {
                coverage = $NF
                gsub(/%/, "", coverage)
                
                if ($0 ~ /\/http\//) { http_total += coverage; http_count++ }
                else if ($0 ~ /\/domain\//) { domain_total += coverage; domain_count++ }
                else if ($0 ~ /\/application\//) { app_total += coverage; app_count++ }
                else if ($0 ~ /\/postgres\//) { repo_total += coverage; repo_count++ }
            }
            END {
                if (http_count > 0) {
                    avg = http_total / http_count
                    if (avg >= 80) printf "🌐 HTTP Layer       : \033[32m%6.1f%%\033[0m (%d functions)\n", avg, http_count
                    else if (avg >= 60) printf "🌐 HTTP Layer       : \033[33m%6.1f%%\033[0m (%d functions)\n", avg, http_count
                    else printf "🌐 HTTP Layer       : \033[31m%6.1f%%\033[0m (%d functions)\n", avg, http_count
                }
                if (domain_count > 0) {
                    avg = domain_total / domain_count
                    if (avg >= 80) printf "🏗️ Domain Layer     : \033[32m%6.1f%%\033[0m (%d functions)\n", avg, domain_count
                    else if (avg >= 60) printf "🏗️ Domain Layer     : \033[33m%6.1f%%\033[0m (%d functions)\n", avg, domain_count
                    else printf "🏗️ Domain Layer     : \033[31m%6.1f%%\033[0m (%d functions)\n", avg, domain_count
                }
                if (app_count > 0) {
                    avg = app_total / app_count
                    if (avg >= 80) printf "🔧 Application Layer: \033[32m%6.1f%%\033[0m (%d functions)\n", avg, app_count
                    else if (avg >= 60) printf "🔧 Application Layer: \033[33m%6.1f%%\033[0m (%d functions)\n", avg, app_count
                    else printf "🔧 Application Layer: \033[31m%6.1f%%\033[0m (%d functions)\n", avg, app_count
                }
                if (repo_count > 0) {
                    avg = repo_total / repo_count
                    if (avg >= 80) printf "💾 Repository Layer : \033[32m%6.1f%%\033[0m (%d functions)\n", avg, repo_count
                    else if (avg >= 60) printf "💾 Repository Layer : \033[33m%6.1f%%\033[0m (%d functions)\n", avg, repo_count
                    else printf "💾 Repository Layer : \033[31m%6.1f%%\033[0m (%d functions)\n", avg, repo_count
                }
            }'
            
            echo ""
            echo "⚠️  Low Coverage Functions (< 60%):"
            echo "─────────────────────────────────────────────────────────────"
            
            # Найдем функции с низким покрытием (только internal/)
            local low_coverage_count=0
            while IFS= read -r line; do
                # Показываем только функции из internal/
                if [[ ! "$line" =~ /internal/ ]]; then
                    continue
                fi
                coverage=$(echo "$line" | awk '{print $NF}' | sed 's/%//')
                if (( $(echo "$coverage < 60" | bc -l) )) && (( $(echo "$coverage > 0" | bc -l) )); then
                    func_name=$(echo "$line" | awk '{print $2}' | sed 's/.*\///')
                    echo -e "\033[31m${func_name}                                    ${coverage}%\033[0m"
                    ((low_coverage_count++))
                    if [ $low_coverage_count -ge 10 ]; then
                        break
                    fi
                fi
            done <<< "$func_coverage"
            
            if [ $low_coverage_count -eq 0 ]; then
                echo "🎉 No functions with low coverage found!"
            fi
            
        else
            echo "⚠️  Could not generate detailed function coverage"
        fi
        
        # Cleanup temporary file
        rm -f "$coverage_file"
    fi
    
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    
    return 0
}

# Функция для полного анализа покрытия
run_full_coverage() {
    echo "🚀 Code Coverage Analysis"
    echo ""
    
    # Анализируем покрытие по слоям
    echo "📋 Running coverage analysis by layers..."
    echo ""
    
    # Domain coverage
    run_coverage_analysis "./tests/domain/..." "Domain Layer Coverage" "domain"
    
    # Handler coverage
    run_coverage_analysis "./tests/integration/cases/quest_handler/..." "Handler Layer Coverage" "handler"
    
    # HTTP coverage
    run_coverage_analysis "./tests/integration/cases/quest_http/..." "HTTP Layer Coverage" "http"
    
    # Общее покрытие
    echo "🎯 Overall Project Coverage:"
    echo "═══════════════════════════════════════════════════════════════"
    go test -p 1 -count=1 -coverprofile=coverage_total.out -coverpkg=./internal/... ./tests/... 2>/dev/null
    
    if [ -f "coverage_total.out" ]; then
        local total_result=$(go tool cover -func=coverage_total.out | tail -1)
        local total_percent=$(echo "$total_result" | awk '{print $NF}')
        # Определяем цвет для общего покрытия проекта
        local project_coverage=$(echo "$total_percent" | sed 's/%//')
        local project_color="\033[31m" # Default to Red
        if (( $(echo "$project_coverage >= 80" | bc -l) )); then
            project_color="\033[32m" # Green for 80%+
        elif (( $(echo "$project_coverage >= 60" | bc -l) )); then
            project_color="\033[33m" # Yellow for 60-79%
        fi
        echo "📊 Total Project Coverage: ${project_color}$total_percent\033[0m"
        
        # Показываем топ непокрытых функций
        echo ""
        echo "🔍 Top Uncovered Functions:"
        echo "─────────────────────────────────────────────────────────────"
        go tool cover -func=coverage_total.out | grep "/internal/" | awk '{
            coverage = $NF
            gsub(/%/, "", coverage)
            if (coverage == 0) {
                func_name = $2
                gsub(/.*\//, "", func_name)
                printf "\033[31m%-50s %6.1f%%\033[0m\n", func_name, coverage
            }
        }' | head -10
        
        rm -f coverage_total.out
    else
        echo "❌ Could not generate total coverage"
    fi
    
    echo ""
    echo "💡 Tips for improving coverage:"
    echo "  • Focus on functions with 0% coverage first"
    echo "  • Add error path testing for command handlers"
    echo "  • Test edge cases in domain validation"
    echo "  • Add repository failure scenarios"
    echo "═══════════════════════════════════════════════════════════════"
}

# Функция для быстрого анализа покрытия
run_quick_coverage() {
    echo "⚡ Quick Coverage Check"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Быстрый анализ общего покрытия
    local output=$(go test -p 1 -coverprofile=quick_coverage.out -coverpkg=./internal/... ./tests/... 2>&1)
    local exit_code=$?
    
    if [ $exit_code -eq 0 ] && [ -f "quick_coverage.out" ]; then
        local total_result=$(go tool cover -func=quick_coverage.out | tail -1)
        local total_percent=$(echo "$total_result" | grep -o "[0-9.]*%")
        
        echo -e "📊 Overall Coverage: \033[32m$total_percent\033[0m"
        
        # Показываем топ-5 лучших и худших
        echo ""
        echo "🏆 Top Covered Functions:"
        go tool cover -func=quick_coverage.out | grep -v "total:" | sort -k3 -nr | head -5 | \
        awk '{printf "✅ %-40s \033[32m%s\033[0m\n", $2, $3}' | sed 's/.*\///'
        
        echo ""
        echo "⚠️  Least Covered Functions:"
        go tool cover -func=quick_coverage.out | grep -v "total:" | awk '{
            coverage = $NF
            gsub(/%/, "", coverage)
            if (coverage > 0 && coverage < 50) print $0
        }' | sort -k3 -n | head -5 | \
        awk '{printf "❌ %-40s \033[31m%s\033[0m\n", $2, $3}' | sed 's/.*\///'
        
        rm -f quick_coverage.out
    else
        echo "❌ Failed to run coverage analysis"
        echo "$output"
    fi
    
    echo "═══════════════════════════════════════════════════════════════"
}

# Обработка аргументов
case "${1:-help}" in
    "full")
        run_full_coverage
        ;;
    "quick")
        run_quick_coverage
        ;;
    "domain")
        run_coverage_analysis "./tests/domain/..." "Domain Layer Coverage" "domain"
        ;;
    "handler")
        run_coverage_analysis "./tests/integration/cases/quest_handler/..." "Handler Layer Coverage" "handler"
        ;;
    "http")
        run_coverage_analysis "./tests/integration/cases/quest_http/..." "HTTP Layer Coverage" "http"
        ;;
    "help"|*)
        echo "🚀 Coverage Analysis Script"
        echo ""
        echo "Usage examples:"
        echo "  ./coverage-report.sh full     # Complete coverage analysis"
        echo "  ./coverage-report.sh quick    # Quick coverage overview"
        echo "  ./coverage-report.sh domain   # Domain layer coverage only"
        echo "  ./coverage-report.sh handler  # Handler layer coverage only"
        echo "  ./coverage-report.sh http     # HTTP layer coverage only"
        echo ""
        ;;
esac