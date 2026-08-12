package config

import "github.com/IBM/sarama"

type PostgresConfig interface {
	URI() string
	MigrationDir() string
}

type OrderHTTPConfig interface {
	Address() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	Config() *sarama.Config
	GroupID() string
}
