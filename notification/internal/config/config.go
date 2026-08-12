package config

import (
	"github.com/joho/godotenv"

	"github.com/H1dEx/ms-rocket/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Kafka                  KafkaConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	Logger                 LoggerConfig
	TelegramBot            TelegramBotConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	kafkaConfig, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}
	orderPaidConsumerConfig, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}
	orderAssembledConsumerConfig, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}
	telegramBotConfig, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}
	appConfig = &config{
		Kafka:                  kafkaConfig,
		OrderPaidConsumer:      orderPaidConsumerConfig,
		OrderAssembledConsumer: orderAssembledConsumerConfig,
		Logger:                 loggerConfig,
		TelegramBot:            telegramBotConfig,
	}
	return nil
}

func GetConfig() *config {
	return appConfig
}
