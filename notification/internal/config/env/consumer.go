package env

import (
	"github.com/IBM/sarama"
)

// newSaramaConsumerConfig creates a standard Sarama consumer configuration
// that is shared across all consumer configs.
func newSaramaConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
