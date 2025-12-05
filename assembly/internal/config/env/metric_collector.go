package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricCollectorEnvConfig struct {
	Endpoint    string        `env:"METRIC_COLLECTOR_ENDPOINT,required"`
	Interval    time.Duration `env:"METRIC_COLLECTOR_INTERVAL,required"`
	ServiceName string        `env:"METRIC_COLLECTOR_SERVICE_NAME,required"`
}

type metricCollectorConfig struct {
	raw metricCollectorEnvConfig
}

func NewMetricCollectorConfig() (*metricCollectorConfig, error) {
	var raw metricCollectorEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricCollectorConfig{raw: raw}, nil
}

func (cfg *metricCollectorConfig) CollectorEndpoint() string {
	return cfg.raw.Endpoint
}

func (cfg *metricCollectorConfig) CollectorInterval() time.Duration {
	return cfg.raw.Interval
}

func (cfg *metricCollectorConfig) ServiceName() string {
	return cfg.raw.ServiceName
}
