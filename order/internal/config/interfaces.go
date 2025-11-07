package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type HTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationsDir() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}
