package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	orderMetrics "github.com/nkolesnikov999/micro2-OK/order/internal/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		start := time.Now()

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		orderMetrics.RequestsTotal.Add(
			r.Context(),
			1,
			metric.WithAttributes(
				attribute.String("method", r.Method),
				attribute.Int("status", rw.statusCode),
			),
		)

		orderMetrics.RequestDuration.Record(
			r.Context(),
			duration,
			metric.WithAttributes(
				attribute.String("method", r.Method),
				attribute.Int("status", rw.statusCode),
			),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
