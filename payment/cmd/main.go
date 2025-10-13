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

const grpcPort = 50050

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

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

// PayOrder обрабатывает запрос на оплату заказа. В реальном сервисе тут должна быть
// интеграция с платёжным провайдером, запись аудита и трассировка
func (s *paymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if req == nil {
		log.Printf("CRITICAL: received nil request in PayOrder - potential infrastructure issue")
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if _, err := uuid.Parse(req.GetOrderUuid()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_uuid format: %v", err)
	}

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

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	service := &paymentService{}
	paymentV1.RegisterPaymentServiceServer(grpcServer, service)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("✅ Server stopped")
}
