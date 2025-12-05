package tracing

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
)

// Константы для заголовков трассировки
const (
	// HTTPTraceIDHeader - заголовок для передачи trace ID в HTTP формате (с большими буквами)
	HTTPTraceIDHeader = "X-Trace-ID"
)

// createHTTPSpanAttributes создает стандартный набор атрибутов для HTTP спана
func createHTTPSpanAttributes(r *http.Request) []trace.SpanStartOption {
	return []trace.SpanStartOption{
		trace.WithAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.URLPath(r.URL.Path),
			semconv.HostName(r.Host),
			semconv.URLScheme(r.URL.Scheme),
			semconv.UserAgentName(r.UserAgent()),
		),
		trace.WithSpanKind(trace.SpanKindServer),
	}
}

// HTTPHandlerMiddleware создает middleware для трассировки HTTP запросов.
// Middleware извлекает контекст трассировки из входящего запроса и создает новый спан для каждого запроса.
// Также добавляет trace ID в заголовки ответа для удобства отладки.
func HTTPHandlerMiddleware(serviceName string) func(http.Handler) http.Handler {
	// Получаем текущий трейсер и пропагатор
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	// Возвращаем функцию-middleware
	return func(next http.Handler) http.Handler {
		// Создаем обработчик HTTP-запросов
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем контекст трассировки из входящего запроса
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Создаем имя операции для спана
			spanName := fmt.Sprintf("%s.%s", r.Method, r.URL.Path)

			// Создаем новый спан для текущего запроса с атрибутами
			ctx, span := tracer.Start(
				ctx,
				spanName,
				createHTTPSpanAttributes(r)...,
			)
			defer span.End()

			// Создаем оболочку для ResponseWriter с добавлением trace ID
			wrw := &traceResponseWriter{
				ResponseWriter: w,
				span:           span,
				headerAdded:    false,
			}

			// Вызываем следующий обработчик с обогащенным контекстом
			next.ServeHTTP(wrw, r.WithContext(ctx))

			// Добавляем статус ответа в атрибуты спана
			span.SetAttributes(semconv.OTelStatusCodeKey.Int(wrw.statusCode))
		})
	}
}
