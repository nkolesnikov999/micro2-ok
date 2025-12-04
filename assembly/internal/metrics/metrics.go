package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "assembly-service"
	namespace   = "micro2-OK"
	appName     = "assembly"
)

// =============================================================================
// METER - ФАБРИКА ДЛЯ СОЗДАНИЯ МЕТРИК
// =============================================================================
//
// Meter в OpenTelemetry - это фабрика для создания инструментов измерения метрик.
// Каждый КОМПОНЕНТ должен иметь свой meter с уникальным именем.
//
// АРХИТЕКТУРА ВЗАИМОДЕЙСТВИЯ:
//
//  1. platform/metrics инициализирует MeterProvider:
//     platform.InitProvider() → otel.SetMeterProvider(meterProvider)
//
//  2. assembly/metrics создает свой Meter:
//     otel.Meter("assembly-service") → получает глобальный MeterProvider
//
//  3. Meter создает метрики через MeterProvider:
//     meter.Int64Counter() → meterProvider.createCounter()
//
//  4. Метрики отправляются через Reader в MeterProvider:
//     Counter.Add() → Reader.collect() → Exporter.export() → OTLP Collector
//
// СХЕМА КОМПОНЕНТОВ:
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                     GLOBAL OTEL REGISTRY                            │
// │  otel.SetMeterProvider(provider) ← platform/metrics                 │
// │  otel.Meter(name) → provider     ← assembly/metrics                 │
// └─────────────────────────────────────────────────────────────────────┘
//
//	↓
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                    METER PROVIDER (один)                            │
// │  ┌─────────────────────┐  ┌─────────────────────┐                   │
// │  │   Reader            │  │   Exporter          │                   │
// │  │ - Периодически      │  │ - Отправляет в      │                   │
// │  │   читает метрики    │  │   OTLP Collector    │                   │
// │  │ - Агрегирует        │  │ - Форматирует       │                   │
// │  │   данные            │  │   протокол          │                   │
// │  └─────────────────────┘  └─────────────────────┘                   │
// └─────────────────────────────────────────────────────────────────────┘
//
//	↓
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                     METERS (много)                                  │
// │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
// │  │ assembly-service │  │ database        │  │ http-client     │      │
// │  │ - RequestsTotal │  │ - Connections   │  │ - Requests      │      │
// │  │ - OrdersTotal   │  │ - QueryDuration │  │ - Errors        │      │
// │  │ - AnalysisTime  │  │ - PoolSize      │  │ - Duration      │      │
// │  └─────────────────┘  └─────────────────┘  └─────────────────┘      │
// └─────────────────────────────────────────────────────────────────────┘
//
// ВАЖНЫЕ ПРИНЦИПЫ:
//
// 1. MeterProvider ОДИН - управляет инфраструктурой отправки метрик
// 2. Meter МНОГО - один на каждый логический компонент (сервис, библиотека)
// 3. Meter получает MeterProvider из глобального registry OpenTelemetry
// 4. Все метрики из всех Meter'ов отправляются через один MeterProvider
// 5. В Prometheus метрики группируются по label'у otel_scope_name
//
// Meter предоставляет методы для создания различных типов метрик:
// - Counter - монотонно возрастающий счетчик
// - UpDownCounter - счетчик, который может увеличиваться и уменьшаться
// - Histogram - распределение значений с bucketing
// - Gauge - моментальное значение (через UpDownCounter или Callback)
//
// Важно: meter должен быть создан один раз и переиспользоваться в рамках компонента
var meter = otel.Meter(serviceName)

// =============================================================================
// ТИПЫ МЕТРИК В OPENTELEMETRY
// =============================================================================
//
// 1. COUNTER (Счетчик) - metric.Int64Counter
//    - Монотонно возрастающее значение (только увеличивается)
//    - Используется для: количество запросов, ошибок, событий
//    - Пример: общее количество HTTP запросов
//    - Методы: Add() - добавить положительное значение
//
// 2. UPDOWNCOUNTER (Двунаправленный счетчик) - metric.Int64UpDownCounter
//    - Может увеличиваться и уменьшаться
//    - Используется для: активные соединения, размер очереди, память
//    - Пример: количество активных gRPC соединений
//    - Методы: Add() - добавить (может быть отрицательным)
//
// 3. HISTOGRAM (Гистограмма) - metric.Float64Histogram
//    - Распределение наблюдений в bucket'ах
//    - Автоматически создает метрики: _count, _sum, _bucket
//    - Используется для: время ответа, размер запроса, задержки
//    - Пример: время выполнения HTTP запроса
//    - Методы: Record() - записать наблюдение
//
// 4. GAUGE (Датчик) - НЕТ отдельного типа в OpenTelemetry!
//    - В OpenTelemetry нет прямого аналога Prometheus Gauge
//    - Для gauge-подобных метрик используются:
//      а) UpDownCounter - когда значение контролируется приложением
//      б) Асинхронные Observable - когда значение нужно читать по требованию
//    - Примеры: температура CPU, использование памяти, размер кэша
//    - Для простых случаев используйте UpDownCounter как gauge

var (
	// MessagesConsumedTotal - COUNTER для подсчета общего количества обработанных сообщений из Kafka
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: подсчет всех сообщений Kafka с разбивкой по топикам и статусам
	// Лейблы: topic (название топика), status (success/error)
	MessagesConsumedTotal metric.Int64Counter

	// MessagesProducedTotal - COUNTER для подсчета общего количества отправленных сообщений в Kafka
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: подсчет всех отправленных сообщений с разбивкой по топикам
	// Лейблы: topic (название топика), status (success/error)
	MessagesProducedTotal metric.Int64Counter

	// MessageProcessingDuration - HISTOGRAM для измерения времени обработки сообщений из Kafka
	// Тип: Float64Histogram (распределение значений)
	// Использование: SLA мониторинг - отслеживание времени обработки сообщений
	// Позволяет строить percentile (p50, p95, p99) для анализа производительности
	MessageProcessingDuration metric.Float64Histogram

	// AssemblyDuration - HISTOGRAM для измерения времени выполнения сборки
	// Тип: Float64Histogram (распределение значений)
	// Использование: отслеживание длительности процесса сборки заказов
	// Bucket boundaries оптимизированы для диапазона 5-20 секунд
	// Автоматически создает метрики: _count, _sum, _bucket для percentile
	AssemblyDuration metric.Float64Histogram
)

// InitMetrics инициализирует все метрики assembly сервиса
// Должна быть вызвана один раз при старте приложения после инициализации OpenTelemetry провайдера
func InitMetrics() error {
	var err error

	// Создаем счетчик обработанных сообщений из Kafka
	MessagesConsumedTotal, err = meter.Int64Counter(
		namespace+"_kafka_"+appName+"_messages_consumed_total",
		metric.WithDescription("Total number of Kafka messages consumed by assembly service"),
	)
	if err != nil {
		return err
	}

	// Создаем счетчик отправленных сообщений в Kafka
	MessagesProducedTotal, err = meter.Int64Counter(
		namespace+"_kafka_"+appName+"_messages_produced_total",
		metric.WithDescription("Total number of Kafka messages produced by assembly service"),
	)
	if err != nil {
		return err
	}

	// Создаем гистограмму времени обработки сообщений из Kafka
	// Bucket'ы оптимизированы для времени обработки сообщений
	MessageProcessingDuration, err = meter.Float64Histogram(
		namespace+"_kafka_"+appName+"_message_processing_duration_seconds",
		metric.WithDescription("Duration of Kafka message processing"),
		metric.WithUnit("s"),
		// Bucket boundaries для обработки сообщений: от миллисекунд до секунд
		// 1ms, 2ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2s, 5s
		metric.WithExplicitBucketBoundaries(
			0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0,
		),
	)
	if err != nil {
		return err
	}

	// Создаем гистограмму времени сборки с bucket'ами для диапазона 5-20 секунд
	AssemblyDuration, err = meter.Float64Histogram(
		namespace+"_"+appName+"_operation_duration_seconds",
		metric.WithDescription("Duration of assembly operations"),
		metric.WithUnit("s"),
		// Bucket boundaries оптимизированы для диапазона 5-20 секунд
		// Позволяет точно измерять распределение времени сборки
		metric.WithExplicitBucketBoundaries(
			5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 12.0, 14.0, 16.0, 18.0, 20.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}
