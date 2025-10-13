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

// PayOrder обрабатывает запрос на оплату заказа и возвращает UUID транзакции.
// В реальном сервисе здесь должна быть интеграция с платёжным провайдером,
// запись аудита и трассировка.
func (s *paymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	// Validate UUID formats
	if _, err := uuid.Parse(req.GetOrderUuid()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_uuid format: %v", err)
	}
	if _, err := uuid.Parse(req.GetUserUuid()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_uuid format: %v", err)
	}

	transactionUuid := uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUuid)

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

	// Регистрируем реализацию сервиса оплаты.
	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)

	// Включаем server reflection для отладки (grpcurl, дебаг).
	reflection.Register(s)

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
