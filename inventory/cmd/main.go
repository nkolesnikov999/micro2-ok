package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partV1API "github.com/nkolesnikov999/micro2-OK/inventory/internal/api/inventory/v1"
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/config"
	partRepository "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/part"
	partService "github.com/nkolesnikov999/micro2-OK/inventory/internal/service/part"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

const configPath = "./deploy/compose/inventory/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v\n", err)
		return
	}
	defer func() {
		cerr := client.Disconnect(ctx)
		if cerr != nil {
			log.Printf("Ошибка отключения от базы данных: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Ошибка проверки соединения с базой данных: %v\n", err)
		return
	}

	db := client.Database("inventory")
	lis, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		log.Printf("Ошибка запуска сервера: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("Ошибка закрытия соединения: %v\n", cerr)
		}
	}()

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	repo := partRepository.NewRepository(ctx, db)
	service := partService.NewService(repo)
	api := partV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(grpcServer, api)

	go func() {
		log.Printf("🚀 gRPC сервер запущен на %s\n", config.AppConfig().GRPC.Address())
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Printf("Ошибка запуска сервера: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Завершение работы gRPC сервера...")
	grpcServer.GracefulStop()
	log.Println("✅ Сервер остановлен")
}
