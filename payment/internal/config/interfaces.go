package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
}

type GRPCConfig interface {
	Address() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
