// Package logger предоставляет dual-write логгер с использованием zapcore.Tee архитектуры
//
// АРХИТЕКТУРА ЛОГГЕРА:
//
// Логгер использует zapcore.NewTee для параллельной записи в два назначения:
// 1. Stdout (для Kubernetes/контейнерных окружений)
// 2. OpenTelemetry коллектор (для централизованного сбора логов)
//
// ПОТОК ДАННЫХ:
//
//		Application
//		    ↓ (logger.Info/Error)
//		zap.Logger
//		    ↓
//		zapcore.Tee
//		   ↙        ↘
//	 StdoutCore   SimpleOTLPCore
//		   ↓             ↓
//	 os.Stdout   SimpleOTLPWriter
//		               ↓
//		        zapcore.BufferedWriteSyncer
//		               ↓
//		         OTLP Collector (gRPC)
//
// КОМПОНЕНТЫ:
//
// 1. StdoutCore - стандартный zap core для вывода в консоль
// 2. SimpleOTLPCore - преобразует zap Entry в OpenTelemetry Record
// 3. SimpleOTLPWriter - отправляет OTLP Records в коллектор
// 4. BufferedWriteSyncer - буферизация для асинхронной отправки
//
// ОСОБЕННОСТИ:
//
// - Graceful degradation: при недоступности OTLP коллектора stdout продолжает работать
// - Метрики: отслеживание sent/dropped записей для мониторинга
// - Батчирование: OTLP SDK автоматически группирует записи для эффективной отправки
// - Таймауты: 500ms лимит для предотвращения блокировки приложения
package logger

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/log"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Глобальные переменные пакета
var (
	global       *zap.Logger                // глобальный экземпляр логгера
	initOnce     sync.Once                  // обеспечивает единократную инициализацию
	level        zap.AtomicLevel            // уровень логирования (может изменяться динамически)
	otelProvider *otelLogSdk.LoggerProvider // OTLP provider для graceful shutdown
)

// Key тип для ключей контекста
type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

// Константы конфигурации OTLP
const (
	serviceEnvironment = "dev" // окружение для фильтрации логов
)

// Таймауты
const (
	shutdownTimeout = 2 * time.Second // таймаут для graceful shutdown OTLP provider
)

// Init инициализирует глобальный логгер с Tee архитектурой.
// Поддерживает одновременную запись в stdout и OTLP коллектор.
//
// Параметры:
//   - logLevel: уровень логирования ("debug", "info", "warn", "error")
//   - asJSON: формат вывода (true - JSON, false - консольный)
//   - enableOTLP: включение отправки в OpenTelemetry коллектор
//   - otlpEndpoint: адрес OTLP коллектора (например, "localhost:4317")
//   - serviceName: имя сервиса для телеметрии
func Init(logLevel string, asJSON, enableOTLP bool, otlpEndpoint, serviceName string) error {
	initOnce.Do(func() {
		level = zap.NewAtomicLevelAt(parseLevel(logLevel))
		cores := buildCores(asJSON, enableOTLP, otlpEndpoint, serviceName)
		global = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	})

	if global == nil {
		return fmt.Errorf("logger init failed")
	}

	return nil
}

// buildCores создает слайс cores для zapcore.Tee.
// Всегда включает stdout core, опционально добавляет OTLP core.
func buildCores(asJSON, enableOTLP bool, otlpEndpoint, serviceName string) []zapcore.Core {
	cores := []zapcore.Core{
		createStdoutCore(asJSON),
	}

	if enableOTLP {
		if otlpCore := createOTLPCore(otlpEndpoint, serviceName); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	return cores
}

// createStdoutCore создает core для записи в stdout/stderr.
// Поддерживает JSON и консольный формат вывода.
func createStdoutCore(asJSON bool) zapcore.Core {
	config := buildEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

// createOTLPCore создает core для отправки в OpenTelemetry коллектор.
// При ошибке подключения возвращает nil (graceful degradation).
func createOTLPCore(otlpEndpoint, serviceName string) *SimpleOTLPCore {
	otlpLogger, err := createOTLPLogger(otlpEndpoint, serviceName)
	if err != nil {
		// Логирование ошибки невозможно, так как логгер еще не инициализирован
		return nil
	}

	// Прямо передаём OTLP-логгер в core. Буферизацию делает OTLP SDK (BatchProcessor).
	return NewSimpleOTLPCore(otlpLogger, level)
}

// createOTLPLogger создает OTLP логгер с настроенным экспортером и ресурсами.
// Использует BatchProcessor для эффективной отправки логов.
func createOTLPLogger(endpoint, serviceName string) (otelLog.Logger, error) {
	ctx := context.Background()

	exporter, err := createOTLPExporter(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	rs, err := createResource(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	provider := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)
	otelProvider = provider // сохраняем для shutdown

	return provider.Logger("app"), nil
}

// createOTLPExporter создает gRPC экспортер для OTLP коллектора
func createOTLPExporter(ctx context.Context, endpoint string) (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(), // для разработки, в продакшене следует использовать TLS
	)
}

// createResource создает метаданные сервиса для телеметрии
func createResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment", serviceEnvironment),
		),
	)
}

// buildEncoderConfig настраивает формат вывода логов с нужными полями
func buildEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}

// Debug записывает лог уровня DEBUG.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Debug(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Debug(msg, fields...)
	}
}

// Info записывает лог уровня INFO.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Info(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Info(msg, fields...)
	}
}

// Warn записывает лог уровня WARN.
func Warn(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Warn(msg, fields...)
	}
}

// Error записывает лог уровня ERROR.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Error(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Error(msg, fields...)
	}
}

// Fatal записывает лог уровня FATAL и завершает процесс.
func Fatal(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Fatal(msg, fields...)
	}
}

// Sync принудительно сбрасывает все буферизованные логи.
// Вызывает sync для всех cores (stdout + OTLP).
func Sync() error {
	if global != nil {
		return global.Sync()
	}

	return nil
}

// Close корректно завершает работу логгера.
// Останавливает OTLP provider с таймаутом для отправки оставшихся логов.
func Close() error {
	if otelProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = otelProvider.Shutdown(ctx)
	}

	return nil
}

// parseLevel преобразует строковое значение в zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// InitForBenchmark инициализирует логгер с no-op core для бенчмарков.
// Используется в тестах производительности, чтобы не засорять вывод и не замедлять тесты.
func InitForBenchmark() {
	initOnce.Do(func() {
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		core := zapcore.NewNopCore() // no-op core для бенчмарков
		global = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})
}

// enrichLogger обёртка над zap.Logger с поддержкой обогащения контекста
type enrichLogger struct {
	zapLogger *zap.Logger
}

// With создает новый enrich-aware логгер с дополнительными полями
func With(fields ...zap.Field) *enrichLogger {
	if global == nil {
		return &enrichLogger{zapLogger: zap.NewNop()}
	}

	return &enrichLogger{
		zapLogger: global.With(fields...),
	}
}

// WithContext создает enrich-aware логгер с полями из контекста
func WithContext(ctx context.Context) *enrichLogger {
	if global == nil {
		return &enrichLogger{zapLogger: zap.NewNop()}
	}

	return &enrichLogger{
		zapLogger: global.With(fieldsFromContext(ctx)...),
	}
}

// Info записывает лог уровня INFO с обогащением из контекста
func (l *enrichLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

// Error записывает лог уровня ERROR с обогащением из контекста
func (l *enrichLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

// fieldsFromContext извлекает enrich-поля из контекста
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}

// logger адаптер, который оборачивает глобальные функции логгера
// для использования в компонентах, требующих интерфейс Logger
type logger struct{}

// Info записывает лог уровня INFO через глобальный логгер
func (g *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	Info(ctx, msg, fields...)
}

// Error записывает лог уровня ERROR через глобальный логгер
func (g *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	Error(ctx, msg, fields...)
}

// Logger возвращает адаптер глобального логгера, реализующий интерфейс Logger
// Используется для передачи логгера в компоненты, которые требуют интерфейс Logger
// (например, closer, kafka consumer/producer, testcontainers и т.д.)
func Logger() *logger {
	return &logger{}
}
