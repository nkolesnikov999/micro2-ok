package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	grpcConn "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/nkolesnikov999/micro2-OK/order/internal/api/order/v1"
	"github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc"
	invClient "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc/inventory/v1"
	payClient "github.com/nkolesnikov999/micro2-OK/order/internal/client/grpc/payment/v1"
	"github.com/nkolesnikov999/micro2-OK/order/internal/config"
	"github.com/nkolesnikov999/micro2-OK/order/internal/repository"
	orderRepository "github.com/nkolesnikov999/micro2-OK/order/internal/repository/order"
	"github.com/nkolesnikov999/micro2-OK/order/internal/service"
	orderService "github.com/nkolesnikov999/micro2-OK/order/internal/service/order"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/closer"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/migrator"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1Server *orderV1.Server

	orderService service.OrderService

	orderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	inventoryConn *grpcConn.ClientConn
	paymentConn   *grpcConn.ClientConn

	postgresDB *pgx.Conn
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1Server(ctx context.Context) (*orderV1.Server, error) {
	if d.orderV1Server == nil {
		orderHandler := orderApi.NewHandler(d.OrderService(ctx))

		server, err := orderV1.NewServer(
			orderHandler,
			orderV1.WithPathPrefix("/api/v1"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create order server: %w", err)
		}
		d.orderV1Server = server
	}

	return d.orderV1Server, nil
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresDB(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) InventoryClient(ctx context.Context) grpc.InventoryClient {
	if d.inventoryClient == nil {
		protoClient := inventoryV1.NewInventoryServiceClient(d.InventoryConn(ctx))
		d.inventoryClient = invClient.NewClient(protoClient)
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) grpc.PaymentClient {
	if d.paymentClient == nil {
		protoClient := paymentV1.NewPaymentServiceClient(d.PaymentConn(ctx))
		d.paymentClient = payClient.NewClient(protoClient)
	}

	return d.paymentClient
}

func (d *diContainer) InventoryConn(ctx context.Context) *grpcConn.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpcConn.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpcConn.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to inventory service: %w", err))
		}

		closer.AddNamed("inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn
}

func (d *diContainer) PaymentConn(ctx context.Context) *grpcConn.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpcConn.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpcConn.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to payment service: %w", err))
		}

		closer.AddNamed("payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn
}

func (d *diContainer) PostgresDB(ctx context.Context) *pgx.Conn {
	if d.postgresDB == nil {
		conn, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to PostgreSQL: %w", err))
		}

		err = conn.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to ping PostgreSQL: %w", err))
		}

		migrationsDir := config.AppConfig().Postgres.MigrationsDir()
		migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*conn.Config().Copy()), migrationsDir)
		err = migratorRunner.Up()
		if err != nil {
			panic(fmt.Errorf("failed to run migrations: %w", err))
		}

		closer.AddNamed("PostgreSQL connection", func(ctx context.Context) error {
			return conn.Close(ctx)
		})

		d.postgresDB = conn
	}

	return d.postgresDB
}
