package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partV1API "github.com/nkolesnikov999/micro2-OK/inventory/internal/api/inventory/v1"
	partRepository "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/part"
	partService "github.com/nkolesnikov999/micro2-OK/inventory/internal/service/part"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Ошибка загрузки .env файла: %v\n", err)
		return
	}

	dbURI := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
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
		log.Printf("🚀 gRPC сервер запущен на порту %d\n", grpcPort)
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
