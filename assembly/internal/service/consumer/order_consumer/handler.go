package order_consumer

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/assembly/internal/config"
	"github.com/nkolesnikov999/micro2-OK/assembly/internal/metrics"
	"github.com/nkolesnikov999/micro2-OK/assembly/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/kafka/consumer"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	// Начало обработки сообщения - измеряем общее время обработки
	processingStart := time.Now()
	consumedTopic := config.AppConfig().OrderPaidConsumer.Topic()

	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		// Метрика: ошибка при обработке сообщения
		metrics.MessagesConsumedTotal.Add(
			ctx,
			1,
			metric.WithAttributes(
				attribute.String("topic", consumedTopic),
				attribute.String("status", "error"),
			),
		)
		processingDuration := time.Since(processingStart).Seconds()
		metrics.MessageProcessingDuration.Record(
			ctx,
			processingDuration,
			metric.WithAttributes(
				attribute.String("topic", consumedTopic),
				attribute.String("status", "error"),
			),
		)
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("payment_method", event.PaymentMethod),
	)

	// start measuring build time from the moment we received the event
	assemblyStart := time.Now()
	// wait random time between 5 and 20 seconds (respecting cancellation)
	//nolint:gosec // G404: math/rand is acceptable for non-cryptographic use case (simulating build time)
	buildDuration := time.Duration(rand.Intn(16)+5) * time.Second // 5-20 seconds
	select {
	case <-time.After(buildDuration):
	case <-ctx.Done():
		// Метрика: обработка прервана
		processingDuration := time.Since(processingStart).Seconds()
		metrics.MessageProcessingDuration.Record(
			ctx,
			processingDuration,
			metric.WithAttributes(
				attribute.String("topic", consumedTopic),
				attribute.String("status", "cancelled"),
			),
		)
		return ctx.Err()
	}

	assemblyElapsed := time.Since(assemblyStart).Seconds()

	// Метрика: время сборки
	metrics.AssemblyDuration.Record(
		ctx,
		assemblyElapsed,
		metric.WithAttributes(
			attribute.String("order_uuid", event.OrderUUID),
		),
	)

	producedTopic := config.AppConfig().OrderAssembledProducer.Topic()
	if err := s.shipAssembledProducer.ProduceShipAssembled(ctx, model.ShipAssembledEvent{
		EventUUID:    uuid.NewString(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(assemblyElapsed),
	}); err != nil {
		logger.Error(ctx, "Failed to produce ShipAssembled", zap.Error(err))
		// Метрика: ошибка при отправке сообщения
		metrics.MessagesProducedTotal.Add(
			ctx,
			1,
			metric.WithAttributes(
				attribute.String("topic", producedTopic),
				attribute.String("status", "error"),
			),
		)
		// Метрика: ошибка при обработке сообщения
		metrics.MessagesConsumedTotal.Add(
			ctx,
			1,
			metric.WithAttributes(
				attribute.String("topic", consumedTopic),
				attribute.String("status", "error"),
			),
		)
		processingDuration := time.Since(processingStart).Seconds()
		metrics.MessageProcessingDuration.Record(
			ctx,
			processingDuration,
			metric.WithAttributes(
				attribute.String("topic", consumedTopic),
				attribute.String("status", "error"),
			),
		)
		return err
	}

	// Метрика: успешно отправлено сообщение в Kafka
	metrics.MessagesProducedTotal.Add(
		ctx,
		1,
		metric.WithAttributes(
			attribute.String("topic", producedTopic),
			attribute.String("status", "success"),
		),
	)

	// Метрика: успешно обработано сообщение
	metrics.MessagesConsumedTotal.Add(
		ctx,
		1,
		metric.WithAttributes(
			attribute.String("topic", consumedTopic),
			attribute.String("status", "success"),
		),
	)

	// Метрика: время обработки сообщения
	processingDuration := time.Since(processingStart).Seconds()
	metrics.MessageProcessingDuration.Record(
		ctx,
		processingDuration,
		metric.WithAttributes(
			attribute.String("topic", consumedTopic),
			attribute.String("status", "success"),
		),
	)

	return nil
}
