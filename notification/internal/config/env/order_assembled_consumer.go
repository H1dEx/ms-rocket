//nolint:dupl // consumer configs differ only by env tags
package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderAssembledConsumerEnvConfig struct {
	Topic   string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type orderAssembledConsumerConfig struct {
	raw orderAssembledConsumerEnvConfig
}

func NewOrderAssembledConsumerConfig() (*orderAssembledConsumerConfig, error) {
	var raw orderAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &orderAssembledConsumerConfig{raw: raw}, nil
}

func (c *orderAssembledConsumerConfig) Topic() string {
	return c.raw.Topic
}

func (c *orderAssembledConsumerConfig) GroupID() string {
	return c.raw.GroupID
}

func (c *orderAssembledConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}
