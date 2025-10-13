// Приложение запускает gRPC‑сервер платежного сервиса.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

// grpcPort — порт, на котором слушает gRPC‑сервер оплаты.
const grpcPort = 50050

// paymentService реализует gRPC‑сервис оплаты заказов.
type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

// isValidPaymentMethod проверяет, является ли способ платежа валидным.
func isValidPaymentMethod(method paymentV1.PaymentMethod) bool {
	switch method {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return true
	default:
		return false
	}
}

// PayOrder обрабатывает запрос на оплату заказа и возвращает UUID транзакции.
// В реальном сервисе здесь должна быть интеграция с платёжным провайдером,
// запись аудита и трассировка.
func (s *paymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	// Check for nil request
	if req == nil {
		log.Printf("CRITICAL: received nil request in PayOrder - potential infrastructure issue")
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// Validate UUID formats
	if _, err := uuid.Parse(req.GetOrderUuid()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_uuid format: %v", err)
	}

	// Validate payment method
	paymentMethod := req.GetPaymentMethod()
	if paymentMethod == paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "payment_method must be specified")
	}
	if !isValidPaymentMethod(paymentMethod) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment_method: %v", paymentMethod)
	}

	transactionUuid := uuid.New().String()
	log.Printf("Оплата прошла успешно: transaction_uuid=%s, order_uuid=%s, user_uuid=%s, payment_method=%s",
		transactionUuid, req.GetOrderUuid(), req.GetUserUuid(), paymentMethod.String())

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}

func main() {
	// Открываем TCP‑порт для gRPC‑сервера.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаём gRPC‑сервер.
	s := grpc.NewServer()

	// Включаем server reflection для отладки (grpcurl, дебаг).
	reflection.Register(s)

	// Регистрируем реализацию сервиса оплаты.
	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)

	// Запускаем сервер в отдельной горутине.
	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Корректное завершение (graceful shutdown): ждём сигнал и останавливаем сервер.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
