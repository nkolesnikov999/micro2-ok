package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/nkolesnikov999/micro2-OK/order/internal/api/order/v1"
	invClient "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc/inventory/v1"
	payClient "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc/payment/v1"
	"github.com/nkolesnikov999/micro2-OK/order/internal/config"
	"github.com/nkolesnikov999/micro2-OK/order/internal/migrator"
	orderRepo "github.com/nkolesnikov999/micro2-OK/order/internal/repository/order"
	orderSvc "github.com/nkolesnikov999/micro2-OK/order/internal/service/order"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

const configPath = "./deploy/compose/order/.env"

func initGRPCConnections() (*grpc.ClientConn, *grpc.ClientConn, error) {
	inventoryConn, err := grpc.NewClient(
		config.AppConfig().InventoryGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось подключиться к inventory сервису: %w", err)
	}

	paymentConn, err := grpc.NewClient(
		config.AppConfig().PaymentGRPC.Address(),
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

func initDatabase(ctx context.Context) (*pgx.Conn, error) {
	con, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
	if err != nil {
		log.Printf("не удалось подключиться к базе данных: %v\n", err)
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	err = con.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return nil, fmt.Errorf("база данных недоступна: %w", err)
	}

	migrationsDir := config.AppConfig().Postgres.MigrationsDir()
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return nil, fmt.Errorf("ошибка миграции базы данных: %w", err)
	}
	return con, nil
}

func initApplication(connDB *pgx.Conn) (*grpc.ClientConn, *grpc.ClientConn, *orderV1.Server, error) {
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
	repo := orderRepo.NewRepository(connDB)
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
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	err = logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}

	ctx := context.Background()
	connDB, err := initDatabase(ctx)
	if err != nil {
		log.Printf("ошибка инициализации базы данных: %v\n", err)
		return
	}
	defer func() {
		cerr := connDB.Close(ctx)
		if cerr != nil {
			log.Printf("ошибка при закрытии соединения с базой данных: %v", cerr)
		}
	}()

	log.Println("🔄 Инициализация приложения...")
	inventoryConn, paymentConn, orderServer, err := initApplication(connDB)
	if err != nil {
		log.Printf("ошибка инициализации приложения: %v", err)
		return
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
		Addr:    config.AppConfig().HTTP.Address(),
		Handler: router,
		// Защита от Slowloris атак: принудительно закрывает соединение, если клиент
		// не успел отправить все заголовки за отведенное время
		ReadHeaderTimeout: config.AppConfig().HTTP.ReadTimeout(),
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на %s\n", config.AppConfig().HTTP.Address())
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig().HTTP.ShutdownTimeout())
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
