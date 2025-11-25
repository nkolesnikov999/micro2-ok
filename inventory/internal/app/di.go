package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	grpcConn "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	partV1API "github.com/nkolesnikov999/micro2-OK/inventory/internal/api/inventory/v1"
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/config"
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/repository"
	partRepository "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/part"
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/service"
	partService "github.com/nkolesnikov999/micro2-OK/inventory/internal/service/part"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/closer"
	grpcAuth "github.com/nkolesnikov999/micro2-OK/platform/pkg/middleware/grpc"
	authV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	partV1API inventoryV1.InventoryServiceServer

	partService service.PartService

	partRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	iamClient       grpcAuth.IAMClient
	iamConn         *grpcConn.ClientConn
	authInterceptor *grpcAuth.AuthInterceptor
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PartV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.partV1API == nil {
		d.partV1API = partV1API.NewAPI(d.PartService(ctx))
	}

	return d.partV1API
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewService(d.PartRepository(ctx))
	}

	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		d.partRepository = partRepository.NewRepository(ctx, d.MongoDBHandle(ctx))
	}

	return d.partRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) IAMConn(ctx context.Context) *grpcConn.ClientConn {
	if d.iamConn == nil {
		conn, err := grpcConn.NewClient(
			config.AppConfig().IAMGRPC.Address(),
			grpcConn.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to IAM service: %w", err))
		}

		closer.AddNamed("IAM gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.iamConn = conn
	}

	return d.iamConn
}

func (d *diContainer) IAMClient(ctx context.Context) grpcAuth.IAMClient {
	if d.iamClient == nil {
		d.iamClient = authV1.NewAuthServiceClient(d.IAMConn(ctx))
	}

	return d.iamClient
}

func (d *diContainer) AuthInterceptor(ctx context.Context) *grpcAuth.AuthInterceptor {
	if d.authInterceptor == nil {
		d.authInterceptor = grpcAuth.NewAuthInterceptor(d.IAMClient(ctx))
	}

	return d.authInterceptor
}
