package config

import "github.com/IBM/sarama"

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	Topic() string
	Config() *sarama.Config
	GroupID() string
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	Config() *sarama.Config
	GroupID() string
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type TelegramBotConfig interface {
	Token() string
}
