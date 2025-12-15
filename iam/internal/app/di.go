package app

import (
	"context"
	"fmt"
	"time"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	authV1API "github.com/nkolesnikov999/micro2-OK/iam/internal/api/auth/v1"
	userV1API "github.com/nkolesnikov999/micro2-OK/iam/internal/api/user/v1"
	"github.com/nkolesnikov999/micro2-OK/iam/internal/config"
	"github.com/nkolesnikov999/micro2-OK/iam/internal/repository"
	sessionRepository "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/session"
	userRepository "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/user"
	"github.com/nkolesnikov999/micro2-OK/iam/internal/service"
	authService "github.com/nkolesnikov999/micro2-OK/iam/internal/service/auth"
	userService "github.com/nkolesnikov999/micro2-OK/iam/internal/service/user"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/cache"
	redisClient "github.com/nkolesnikov999/micro2-OK/platform/pkg/cache/redis"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/closer"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/migrator"
	authV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
	userV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authV1API authV1.AuthServiceServer
	userV1API userV1.UserServiceServer

	authService service.AuthService
	userService service.UserService

	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository

	postgresDB  *pgx.Conn
	redisClient cache.RedisClient
	redisPool   *redigo.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authV1API.NewAPI(d.AuthService(ctx))
	}

	return d.authV1API
}

// AuthorizationServer возвращает AuthorizationServer для ext_authz (Envoy External Authorization)
// Тот же объект, что и AuthV1API, так как api структура реализует оба интерфейса
func (d *diContainer) AuthorizationServer(ctx context.Context) authv3.AuthorizationServer {
	// Используем type assertion, так как api структура реализует оба интерфейса
	return d.AuthV1API(ctx).(authv3.AuthorizationServer)
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userV1API.NewAPI(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(
			d.SessionRepository(ctx),
			d.UserRepository(ctx),
		)
	}

	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(
			d.UserRepository(ctx),
		)
	}

	return d.userService
}

func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisClient(ctx))
	}

	return d.sessionRepository
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PostgresDB(ctx))
	}

	return d.userRepository
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

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		redisCfg := config.AppConfig().Redis
		d.redisPool = &redigo.Pool{
			MaxIdle:     redisCfg.MaxIdle(),
			IdleTimeout: redisCfg.IdleTimeout(),
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", redisCfg.Address())
			},
			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}

		closer.AddNamed("Redis pool", func(ctx context.Context) error {
			return d.redisPool.Close()
		})
	}

	return d.redisPool
}

func (d *diContainer) RedisClient(ctx context.Context) cache.RedisClient {
	if d.redisClient == nil {
		redisCfg := config.AppConfig().Redis
		d.redisClient = redisClient.NewClient(
			d.RedisPool(),
			logger.Logger(),
			redisCfg.ConnectionTimeout(),
		)
	}

	return d.redisClient
}
