package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/nkolesnikov999/micro2-OK/order/internal/api/order/v1"
	invClient "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc/inventory/v1"
	payClient "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/nkolesnikov999/micro2-OK/order/internal/repository/order"
	orderSvc "github.com/nkolesnikov999/micro2-OK/order/internal/service/order"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	inventoryAddr     = "127.0.0.1:50051"
	paymentAddr       = "127.0.0.1:50050"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func initGRPCConnections() (*grpc.ClientConn, *grpc.ClientConn, error) {
	inventoryConn, err := grpc.NewClient(
		inventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось подключиться к inventory сервису: %w", err)
	}

	paymentConn, err := grpc.NewClient(
		paymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		// Cleanup: закрываем уже открытое соединение при ошибке
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с inventory: %v", cerr)
		}
		return nil, nil, fmt.Errorf("не удалось подключиться к payment сервису: %w", err)
	}

	return inventoryConn, paymentConn, nil
}

func initApplication() (*grpc.ClientConn, *grpc.ClientConn, *orderV1.Server, error) {
	inventoryConn, paymentConn, err := initGRPCConnections()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("ошибка инициализации gRPC соединений: %w", err)
	}

	// External gRPC clients
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	// Internal adapters
	grpcInventory := invClient.NewClient(inventoryClient)
	grpcPayment := payClient.NewClient(paymentClient)

	// Repository, Service, API handler
	repo := orderRepo.NewRepository()
	svc := orderSvc.NewService(repo, grpcInventory, grpcPayment)
	handler := orderApi.NewHandler(svc)

	orderServer, err := orderV1.NewServer(
		handler,
		orderV1.WithPathPrefix("/api/v1"),
	)
	if err != nil {
		// Cleanup: закрываем соединения при ошибке создания сервера
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с inventory: %v", cerr)
		}
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с payment: %v", cerr)
		}
		return nil, nil, nil, fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	return inventoryConn, paymentConn, orderServer, nil
}

func main() {
	inventoryConn, paymentConn, orderServer, err := initApplication()
	if err != nil {
		log.Fatalf("ошибка инициализации приложения: %v", err)
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с inventory: %v", cerr)
		}
	}()
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с payment: %v", cerr)
		}
	}()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Mount("/", orderServer)

	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", httpPort),
		Handler: router,
		// Защита от Slowloris атак: принудительно закрывает соединение, если клиент
		// не успел отправить все заголовки за отведенное время
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
