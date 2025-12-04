package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "order-service"
	namespace   = "micro2-OK"
	appName     = "order"
)

var meter = otel.Meter(serviceName)

var (
	// RequestsTotal - COUNTER для подсчета общего количества HTTP запросов
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: подсчет всех HTTP запросов с разбивкой по методам и статусам
	// Лейблы: method (HTTP метод: GET, POST, PUT, DELETE), status (HTTP статус код)
	RequestsTotal metric.Int64Counter

	// OrdersTotal - COUNTER для подсчета созданных заказов
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: бизнес-метрика для отслеживания количества новых заказов
	// Лейблы: нет (простой счетчик без группировки)
	OrdersTotal metric.Int64Counter

	// OrdersRevenueTotal - COUNTER для подсчета суммарной выручки
	// Тип: Float64Counter (монотонно возрастающий)
	// Использование: бизнес-метрика для отслеживания общей выручки от всех заказов
	// Лейблы: нет (простой счетчик без группировки)
	OrdersRevenueTotal metric.Float64Counter

	// AnalysisRequestsTotal - COUNTER для подсчета запросов на анализ
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: отслеживание использования функции анализа заказов
	// Лейблы: нет (простой счетчик)
	AnalysisRequestsTotal metric.Int64Counter

	// AnalysisDuration - HISTOGRAM для измерения времени выполнения анализа
	// Тип: Float64Histogram (распределение значений)
	// Использование: отслеживание производительности алгоритма анализа
	// Автоматически создает метрики: _count, _sum, _bucket для percentile
	AnalysisDuration metric.Float64Histogram

	// RequestDuration - HISTOGRAM для измерения времени выполнения HTTP запросов
	// Тип: Float64Histogram (распределение значений)
	// Использование: SLA мониторинг - отслеживание времени ответа HTTP API
	// Позволяет строить percentile (p50, p95, p99) для анализа производительности
	RequestDuration metric.Float64Histogram
)

// InitMetrics инициализирует все метрики order сервиса
// Должна быть вызвана один раз при старте приложения после инициализации OpenTelemetry провайдера
func InitMetrics() error {
	var err error

	// Создаем счетчик запросов с описанием для документации
	RequestsTotal, err = meter.Int64Counter(
		namespace+"_http_"+appName+"_requests_total",
		metric.WithDescription("Total number of order service requests"),
	)
	if err != nil {
		return err
	}

	// Создаем счетчик заказов
	OrdersTotal, err = meter.Int64Counter(
		namespace+"_"+appName+"_orders_total",
		metric.WithDescription("Total number of orders created"),
	)
	if err != nil {
		return err
	}

	// Создаем счетчик суммарной выручки
	OrdersRevenueTotal, err = meter.Float64Counter(
		namespace+"_"+appName+"_orders_revenue_total",
		metric.WithDescription("Total revenue from all orders"),
		metric.WithUnit("currency"),
	)
	if err != nil {
		return err
	}

	// Создаем гистограмму времени HTTP запросов с правильными bucket'ами
	// Bucket'ы оптимизированы для времени отклика HTTP в диапазоне от миллисекунд до секунд
	RequestDuration, err = meter.Float64Histogram(
		namespace+"_http_"+appName+"_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("s"),
		// Добавляем explicit bucket boundaries для более точного измерения HTTP запросов
		// 1ms, 2ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2s, 5s
		metric.WithExplicitBucketBoundaries(
			0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}
