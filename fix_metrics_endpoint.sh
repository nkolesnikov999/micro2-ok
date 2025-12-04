#!/bin/bash

# Скрипт для исправления endpoint метрик в .env файлах
# Меняет otel-collector на localhost для приложений, запущенных вне контейнера

echo "Исправление METRIC_COLLECTOR_ENDPOINT в .env файлах..."
echo ""

FILES=(
    "deploy/compose/order/.env"
    "deploy/compose/assembly/.env"
)

for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "Обработка $file..."
        # Создаем backup
        cp "$file" "${file}.bak"
        
        # Заменяем otel-collector на localhost
        sed -i 's|METRIC_COLLECTOR_ENDPOINT=http://otel-collector:4318|METRIC_COLLECTOR_ENDPOINT=http://localhost:4318|g' "$file"
        
        echo "✅ Исправлено: $file"
        echo "   Backup сохранен: ${file}.bak"
    else
        echo "⚠️  Файл не найден: $file"
    fi
done

echo ""
echo "Готово! Теперь перезапустите приложения order и assembly."
echo "После перезапуска подождите 10-15 секунд и запустите: ./check_metrics.sh"

