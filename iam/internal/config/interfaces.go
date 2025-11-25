package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationsDir() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type SessionConfig interface {
	TTL() time.Duration
}
