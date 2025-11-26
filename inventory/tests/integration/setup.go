//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers/app"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers/mongo"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers/network"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers/path"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers/postgres"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers/redis"
)

const (
	// Параметры для контейнеров
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

// TestEnvironment type is defined in test_environment.go

// setupTestEnvironment — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	// Получаем переменные окружения для PostgreSQL (с дефолтными значениями)
	postgresImageName := getEnvWithDefault(ctx, testcontainers.PostgresImageNameKey, "postgres:15")
	postgresUser := getEnvWithDefault(ctx, testcontainers.PostgresUserKey, "postgres")
	postgresPassword := getEnvWithDefault(ctx, testcontainers.PostgresPasswordKey, "postgres")
	postgresDatabase := getEnvWithDefault(ctx, testcontainers.PostgresDatabaseKey, "iam_db")

	// Получаем переменные окружения для Redis (с дефолтным значением)
	redisImageName := getEnvWithDefault(ctx, testcontainers.RedisImageNameKey, "redis:7-alpine")

	// Получаем порт gRPC для waitStrategy
	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	// Шаг 2: Запускаем контейнер с MongoDB
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB успешно запущен")

	// Шаг 2.1: Запускаем контейнер с PostgreSQL
	postgresOpts := []postgres.Option{
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithContainerName(testcontainers.PostgresContainerName),
		postgres.WithDatabase(postgresDatabase),
		postgres.WithAuth(postgresUser, postgresPassword),
		postgres.WithLogger(logger.Logger()),
	}
	if postgresImageName != "" {
		postgresOpts = append(postgresOpts, postgres.WithImageName(postgresImageName))
	}
	generatedPostgres, err := postgres.NewContainer(ctx, postgresOpts...)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось запустить контейнер PostgreSQL", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер PostgreSQL успешно запущен")

	// Шаг 2.2: Запускаем контейнер с Redis
	redisOpts := []redis.Option{
		redis.WithNetworkName(generatedNetwork.Name()),
		redis.WithContainerName(testcontainers.RedisContainerName),
		redis.WithLogger(logger.Logger()),
	}
	if redisImageName != "" {
		redisOpts = append(redisOpts, redis.WithImageName(redisImageName))
	}
	generatedRedis, err := redis.NewContainer(ctx, redisOpts...)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось запустить контейнер Redis", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер Redis успешно запущен")

	// Шаг 3: Запускаем контейнер с IAM приложением
	projectRoot := path.GetProjectRoot()

	iamEnv := map[string]string{
		// PostgreSQL настройки
		"POSTGRES_HOST": generatedPostgres.Config().ContainerName,
		// Используем внутренний порт контейнера PostgreSQL, доступный внутри Docker-сети
		"POSTGRES_PORT":       testcontainers.PostgresPort,
		"POSTGRES_DB":         postgresDatabase,
		"POSTGRES_USER":       postgresUser,
		"POSTGRES_PASSWORD":   postgresPassword,
		"POSTGRES_SSL_MODE":   "disable",
		"MIGRATION_DIRECTORY": "./iam/migrations",
		// Redis настройки
		"REDIS_HOST": generatedRedis.Config().ContainerName,
		// Используем внутренний порт контейнера Redis, доступный внутри Docker-сети
		"REDIS_PORT":               testcontainers.RedisPort,
		"REDIS_CONNECTION_TIMEOUT": "5s",
		"REDIS_MAX_IDLE":           "10",
		"REDIS_IDLE_TIMEOUT":       "5m",
		"REDIS_CACHE_TTL":          "1h",
		// gRPC настройки
		"GRPC_HOST": "0.0.0.0",
		"GRPC_PORT": iamGRPCPort,
		// Session настройки
		"SESSION_TTL": "24h",
		// Logger настройки
		"LOGGER_LEVEL":   "debug",
		"LOGGER_AS_JSON": "true",
	}

	iamWaitStrategy := wait.ForListeningPort(nat.Port(iamGRPCPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	iamContainer, err := app.NewContainer(ctx,
		app.WithName(iamAppName),
		app.WithPort(iamGRPCPort),
		app.WithDockerfile(projectRoot, iamDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(iamEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(iamWaitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{
			Network:  generatedNetwork,
			Mongo:    generatedMongo,
			Postgres: generatedPostgres,
			Redis:    generatedRedis,
		})
		logger.Fatal(ctx, "не удалось запустить контейнер IAM", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер IAM успешно запущен")

	// Шаг 4: Запускаем контейнер с Inventory приложением
	appEnv := map[string]string{
		// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
		// IAM gRPC настройки - используем имя контейнера для внутренней сети
		"IAM_GRPC_HOST": iamAppName,
		"IAM_GRPC_PORT": iamGRPCPort,
		// Logger настройки для inventory
		"LOGGER_LEVEL":   "debug",
		"LOGGER_AS_JSON": "true",
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{
			Network:  generatedNetwork,
			Mongo:    generatedMongo,
			Postgres: generatedPostgres,
			Redis:    generatedRedis,
			IAM:      iamContainer,
		})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network:  generatedNetwork,
		Mongo:    generatedMongo,
		App:      appContainer,
		IAM:      iamContainer,
		Postgres: generatedPostgres,
		Redis:    generatedRedis,
	}
}

// getEnvWithLogging возвращает значение переменной окружения с логированием
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}

// getEnvWithDefault возвращает значение переменной окружения или дефолтное значение
func getEnvWithDefault(ctx context.Context, key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
