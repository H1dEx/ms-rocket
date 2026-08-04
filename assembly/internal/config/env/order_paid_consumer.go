package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidConsumerEnvConfig struct {
	Topic   string `env:"ASSEMBLY_ORDER_PAID_TOPIC_NAME"`
	GroupID string `env:"ASSEMBLY_ORDER_PAID_CONSUMER_GROUP_ID"`
}

type orderPaidConsumerConfig struct {
	raw orderPaidConsumerEnvConfig
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidConsumerConfig{
		raw: raw,
	}, nil
}

func (c *orderPaidConsumerConfig) Topic() string {
	return c.raw.Topic
}

func (c *orderPaidConsumerConfig) GroupID() string {
	return c.raw.GroupID
}

func (c *orderPaidConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
