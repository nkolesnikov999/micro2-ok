package redis

import (
	"context"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startRedisContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:               cfg.ContainerName,
		Image:              cfg.ImageName,
		Networks:           []string{cfg.NetworkName},
		WaitingFor:         wait.ForListeningPort(redisPort + "/tcp").WithStartupTimeout(redisStartupTimeout),
		HostConfigModifier: defaultHostConfig(),
		ExposedPorts:       []string{redisPort + "/tcp"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Errorf("failed to start redis container: %v", err)
	}

	return container, nil
}

func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", errors.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, redisPort+"/tcp")
	if err != nil {
		return "", "", errors.Errorf("failed to get mapped port: %v", err)
	}

	return host, port.Port(), nil
}
