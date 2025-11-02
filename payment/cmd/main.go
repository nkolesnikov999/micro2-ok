package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1API "github.com/nkolesnikov999/micro2-OK/payment/internal/api/payment/v1"
	"github.com/nkolesnikov999/micro2-OK/payment/internal/config"
	paymentService "github.com/nkolesnikov999/micro2-OK/payment/internal/service/payment"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/grpc/health"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

const configPath = "./deploy/compose/payment/.env"

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

	lis, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
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

	health.RegisterService(grpcServer)

	service := paymentService.NewService()
	api := paymentV1API.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(grpcServer, api)

	go func() {
		log.Printf("🚀 gRPC server listening on %s\n", config.AppConfig().GRPC.Address())
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
