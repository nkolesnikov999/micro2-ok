package tracing

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// traceResponseWriter - это обертка для http.ResponseWriter, которая отслеживает статус ответа
// и добавляет trace ID в заголовки ответа.
type traceResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	span        trace.Span
	headerAdded bool
}

// addTraceIDHeader добавляет trace ID в заголовки ответа, если он ещё не был добавлен
func (w *traceResponseWriter) addTraceIDHeader() {
	if !w.headerAdded {
		traceID := w.span.SpanContext().TraceID().String()
		if traceID != "" {
			w.ResponseWriter.Header().Set(HTTPTraceIDHeader, traceID)
		}
		w.headerAdded = true
	}
}

// WriteHeader перехватывает запись заголовка ответа для сохранения статус-кода
// и добавляет trace ID в заголовки.
func (w *traceResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.addTraceIDHeader()
	w.ResponseWriter.WriteHeader(code)
}

// Write перехватывает запись тела ответа и устанавливает статус 200, если заголовок еще не был записан.
func (w *traceResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	w.addTraceIDHeader()
	return w.ResponseWriter.Write(b)
}
