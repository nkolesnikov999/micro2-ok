#!/bin/bash

# Скрипт для проверки наличия метрик в Prometheus
# Проверяет все метрики из order и assembly сервисов

PROMETHEUS_URL="${PROMETHEUS_URL:-http://localhost:9090}"

echo "=========================================="
echo "Проверка метрик в Prometheus"
echo "=========================================="
echo "Prometheus URL: $PROMETHEUS_URL"
echo ""

# Проверка доступности Prometheus
echo "1. Проверка доступности Prometheus..."
HEALTH=$(curl -s "$PROMETHEUS_URL/-/healthy")
if [ $? -ne 0 ]; then
    echo "❌ Prometheus недоступен!"
    exit 1
fi
echo "✅ Prometheus доступен: $HEALTH"
echo ""

# Проверка статуса OpenTelemetry Collector
echo "2. Проверка статуса OpenTelemetry Collector..."
OTEL_COLLECTOR_STATUS=$(curl -s "$PROMETHEUS_URL/api/v1/query?query=up{job=\"otel-collector\"}" 2>/dev/null | python3 -c "import sys, json; data = json.load(sys.stdin); results = data.get('data', {}).get('result', []); print('1' if any(r.get('value', [None, '0'])[1] == '1' for r in results) else '0')" 2>/dev/null || echo "0")
if [ "$OTEL_COLLECTOR_STATUS" = "1" ]; then
    echo "✅ OpenTelemetry Collector работает"
else
    echo "⚠️  OpenTelemetry Collector недоступен (up=0)"
    echo "   Это может быть причиной отсутствия метрик!"
fi
echo ""

# Получение списка всех метрик
echo "3. Получение списка всех метрик..."
ALL_METRICS=$(curl -s "$PROMETHEUS_URL/api/v1/label/__name__/values")
if [ $? -ne 0 ]; then
    echo "❌ Не удалось получить список метрик!"
    exit 1
fi
echo "✅ Список метрик получен"
echo ""

# Список ожидаемых метрик из order сервиса
ORDER_METRICS=(
    "micro2-OK_http_order_requests_total"
    "micro2-OK_order_orders_total"
    "micro2-OK_order_orders_revenue_total"
    "micro2-OK_http_order_request_duration_seconds_count"
    "micro2-OK_http_order_request_duration_seconds_sum"
    "micro2-OK_http_order_request_duration_seconds_bucket"
)

# Список ожидаемых метрик из assembly сервиса
ASSEMBLY_METRICS=(
    "micro2-OK_kafka_assembly_messages_consumed_total"
    "micro2-OK_kafka_assembly_messages_produced_total"
    "micro2-OK_kafka_assembly_message_processing_duration_seconds_count"
    "micro2-OK_kafka_assembly_message_processing_duration_seconds_sum"
    "micro2-OK_kafka_assembly_message_processing_duration_seconds_bucket"
    "micro2-OK_assembly_operation_duration_seconds_count"
    "micro2-OK_assembly_operation_duration_seconds_sum"
    "micro2-OK_assembly_operation_duration_seconds_bucket"
)

# Функция для проверки наличия метрики
check_metric() {
    local metric_name=$1
    if echo "$ALL_METRICS" | grep -q "\"$metric_name\""; then
        echo "  ✅ $metric_name"
        return 0
    else
        echo "  ❌ $metric_name - НЕ НАЙДЕНА"
        return 1
    fi
}

# Проверка метрик order сервиса
echo "4. Проверка метрик Order сервиса:"
ORDER_FOUND=0
ORDER_TOTAL=${#ORDER_METRICS[@]}
for metric in "${ORDER_METRICS[@]}"; do
    if check_metric "$metric"; then
        ((ORDER_FOUND++))
    fi
done
echo "   Найдено: $ORDER_FOUND из $ORDER_TOTAL"
echo ""

# Проверка метрик assembly сервиса
echo "5. Проверка метрик Assembly сервиса:"
ASSEMBLY_FOUND=0
ASSEMBLY_TOTAL=${#ASSEMBLY_METRICS[@]}
for metric in "${ASSEMBLY_METRICS[@]}"; do
    if check_metric "$metric"; then
        ((ASSEMBLY_FOUND++))
    fi
done
echo "   Найдено: $ASSEMBLY_FOUND из $ASSEMBLY_TOTAL"
echo ""

# Итоговая статистика
TOTAL_FOUND=$((ORDER_FOUND + ASSEMBLY_FOUND))
TOTAL_EXPECTED=$((ORDER_TOTAL + ASSEMBLY_TOTAL))

echo "=========================================="
echo "Итоговая статистика:"
echo "=========================================="
echo "Order сервис:    $ORDER_FOUND/$ORDER_TOTAL метрик"
echo "Assembly сервис: $ASSEMBLY_FOUND/$ASSEMBLY_TOTAL метрик"
echo "Всего:           $TOTAL_FOUND/$TOTAL_EXPECTED метрик"
echo ""

if [ $TOTAL_FOUND -eq $TOTAL_EXPECTED ]; then
    echo "✅ Все метрики найдены в Prometheus!"
    exit 0
else
    echo "⚠️  Не все метрики найдены в Prometheus"
    echo ""
    echo "=========================================="
    echo "Диагностика проблемы:"
    echo "=========================================="
    
    # Проверка наличия других метрик приложения
    echo "6. Поиск похожих метрик в Prometheus..."
    SIMILAR_METRICS=$(echo "$ALL_METRICS" | python3 -c "import sys, json; data = json.load(sys.stdin); metrics = [m for m in data.get('data', []) if any(x in m.lower() for x in ['order', 'assembly', 'kafka', 'http', 'request', 'message'])]; print('\n'.join(metrics[:10]))" 2>/dev/null)
    if [ -n "$SIMILAR_METRICS" ]; then
        echo "Найдены похожие метрики (возможно, с другими именами):"
        echo "$SIMILAR_METRICS" | sed 's/^/  - /'
    else
        echo "  Похожие метрики не найдены"
    fi
    echo ""
    
    echo "Возможные причины:"
    echo "  1. ⚠️  OpenTelemetry Collector не запущен или недоступен"
    echo "     Решение: запустите 'task up-core' для запуска Collector"
    echo ""
    echo "  2. Сервисы не отправляют метрики в Collector"
    echo "     Проверьте:"
    echo "     - Запущены ли сервисы order и assembly"
    echo "     - Правильно ли настроен METRIC_COLLECTOR_ENDPOINT"
    echo "     - Есть ли ошибки в логах сервисов"
    echo ""
    echo "  3. Метрики еще не были отправлены"
    echo "     Интервал отправки обычно 5-10 секунд"
    echo "     Подождите и запустите скрипт снова"
    echo ""
    echo "  4. Проблемы с конфигурацией экспорта"
    echo "     Проверьте collector.yaml и prometheus.yml"
    echo ""
    echo "Для запуска OpenTelemetry Collector выполните:"
    echo "  task up-core"
    echo ""
    echo "Или вручную:"
    echo "  cd deploy/compose/core && docker compose up -d otel-collector"
    
    exit 1
fi

