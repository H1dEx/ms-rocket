package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	"github.com/H1dEx/ms-rocket/notification/internal/client/http"
	"github.com/H1dEx/ms-rocket/notification/internal/client/http/telegram"
	"github.com/H1dEx/ms-rocket/notification/internal/config"
	decoderDef "github.com/H1dEx/ms-rocket/notification/internal/converter/kafka"
	decoder "github.com/H1dEx/ms-rocket/notification/internal/converter/kafka/decoder"
	"github.com/H1dEx/ms-rocket/notification/internal/service"
	assembledConsumer "github.com/H1dEx/ms-rocket/notification/internal/service/consumer/order_assembled_consumer"
	paidConsumer "github.com/H1dEx/ms-rocket/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/H1dEx/ms-rocket/notification/internal/service/telegram"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka/consumer"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

type diContainer struct {
	orderPaidConsumer      service.OrderPaidService
	orderAssembledConsumer service.OrderAssembledService
	telegramService        service.TelegramService

	paidDecoder            decoderDef.OrderPaidDecoder
	assembledDecoder       decoderDef.OrderAssembledDecoder
	kafkaPaidConsumer      kafka.Consumer
	kafkaAssembledConsumer kafka.Consumer

	paidConsumerGroup      sarama.ConsumerGroup
	assembledConsumerGroup sarama.ConsumerGroup

	telegramClient http.TelegramClient
	telegramBot    *bot.Bot
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) OrderPaidConsumer(ctx context.Context) service.OrderPaidService {
	if c.orderPaidConsumer == nil {
		c.orderPaidConsumer = paidConsumer.NewService(c.KafkaPaidConsumer(ctx), c.PaidDecoder(ctx), c.TelegramService(ctx))
	}
	return c.orderPaidConsumer
}

func (c *diContainer) PaidDecoder(_ context.Context) decoderDef.OrderPaidDecoder {
	if c.paidDecoder == nil {
		c.paidDecoder = decoder.NewOrderPaidDecoder()
	}
	return c.paidDecoder
}

func (c *diContainer) KafkaPaidConsumer(ctx context.Context) kafka.Consumer {
	if c.kafkaPaidConsumer == nil {
		c.kafkaPaidConsumer = consumer.NewConsumer(c.PaidConsumerGroup(ctx), []string{config.GetConfig().OrderPaidConsumer.Topic()}, logger.Logger())
	}
	return c.kafkaPaidConsumer
}

func (c *diContainer) PaidConsumerGroup(_ context.Context) sarama.ConsumerGroup {
	if c.paidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(config.GetConfig().Kafka.Brokers(), config.GetConfig().OrderPaidConsumer.GroupID(), config.GetConfig().OrderPaidConsumer.Config())
		if err != nil {
			panic(fmt.Errorf("failed to create paid consumer group: %s", err.Error()))
		}
		closer.AddNamed("paid_consumer_group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})
		c.paidConsumerGroup = consumerGroup
	}
	return c.paidConsumerGroup
}

func (c *diContainer) OrderAssembledConsumer(ctx context.Context) service.OrderAssembledService {
	if c.orderAssembledConsumer == nil {
		c.orderAssembledConsumer = assembledConsumer.NewService(c.KafkaAssembledConsumer(ctx), c.AssembledDecoder(ctx), c.TelegramService(ctx))
	}
	return c.orderAssembledConsumer
}

func (c *diContainer) AssembledDecoder(_ context.Context) decoderDef.OrderAssembledDecoder {
	if c.assembledDecoder == nil {
		c.assembledDecoder = decoder.NewOrderAssembledDecoder()
	}
	return c.assembledDecoder
}

func (c *diContainer) KafkaAssembledConsumer(ctx context.Context) kafka.Consumer {
	if c.kafkaAssembledConsumer == nil {
		c.kafkaAssembledConsumer = consumer.NewConsumer(c.AssembledConsumerGroup(ctx), []string{config.GetConfig().OrderAssembledConsumer.Topic()}, logger.Logger())
	}
	return c.kafkaAssembledConsumer
}

func (c *diContainer) AssembledConsumerGroup(_ context.Context) sarama.ConsumerGroup {
	if c.assembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(config.GetConfig().Kafka.Brokers(), config.GetConfig().OrderAssembledConsumer.GroupID(), config.GetConfig().OrderAssembledConsumer.Config())
		if err != nil {
			panic(fmt.Errorf("failed to create assembled consumer group: %s", err.Error()))
		}
		closer.AddNamed("assembled_consumer_group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})
		c.assembledConsumerGroup = consumerGroup
	}
	return c.assembledConsumerGroup
}

func (c *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if c.telegramService == nil {
		c.telegramService = telegramService.NewService(c.TelegramClient(ctx))
	}
	return c.telegramService
}

func (c *diContainer) TelegramClient(ctx context.Context) http.TelegramClient {
	if c.telegramClient == nil {
		c.telegramClient = telegram.NewClient(c.TelegramBot(ctx))
	}
	return c.telegramClient
}

func (c *diContainer) TelegramBot(_ context.Context) *bot.Bot {
	if c.telegramBot == nil {
		bot, err := bot.New(config.GetConfig().TelegramBot.Token())
		if err != nil {
			panic(fmt.Errorf("failed to create telegram bot: %s", err.Error()))
		}
		c.telegramBot = bot
	}
	return c.telegramBot
}
