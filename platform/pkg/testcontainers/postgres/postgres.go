package postgres

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

const (
	postgresPort           = "5432"
	postgresStartupTimeout = 1 * time.Minute

	postgresEnvUserKey     = "POSTGRES_USER"
	postgresEnvPasswordKey = "POSTGRES_PASSWORD" //nolint:gosec
	postgresEnvDatabaseKey = "POSTGRES_DB"
)

type Container struct {
	container testcontainers.Container
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := startPostgresContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Postgres container started", zap.String("host", cfg.Host), zap.String("port", cfg.Port))
	success = true

	return &Container{
		container: container,
		cfg:       cfg,
	}, nil
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) Terminate(ctx context.Context) error {
	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
		return err
	}

	c.cfg.Logger.Info(ctx, "Postgres container terminated")

	return nil
}
