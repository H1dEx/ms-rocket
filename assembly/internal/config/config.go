package config

import (
	"github.com/joho/godotenv"

	"github.com/H1dEx/ms-rocket/assembly/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderAssembledProducer OrderAssembledProducerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}
	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}
	kafkaConfig, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderAssembledProducerConfig, err := env.NewOrderAssembledProducerConfig()
	if err != nil {
		return err
	}
	orderPaidConsumerConfig, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerConfig,
		Kafka:                  kafkaConfig,
		OrderAssembledProducer: orderAssembledProducerConfig,
		OrderAssembledConsumer: orderPaidConsumerConfig,
	}
	return nil
}

func GetConfig() *config {
	return appConfig
}
