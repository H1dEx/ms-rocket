package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/H1dEx/ms-rocket/assembly/internal/config"
	kafka_decoder "github.com/H1dEx/ms-rocket/assembly/internal/converter/kafka"
	"github.com/H1dEx/ms-rocket/assembly/internal/converter/kafka/decoder"
	"github.com/H1dEx/ms-rocket/assembly/internal/service"
	"github.com/H1dEx/ms-rocket/assembly/internal/service/consumer/order_consumer"
	"github.com/H1dEx/ms-rocket/assembly/internal/service/producer/order_producer"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	platform_kafka "github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	platform_kafka_consumer "github.com/H1dEx/ms-rocket/platform/pkg/kafka/consumer"
	platform_kafka_producer "github.com/H1dEx/ms-rocket/platform/pkg/kafka/producer"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

type diContainer struct {
	consumerService service.ConsumerService
	producerService service.ProducerService

	consumerGroup sarama.ConsumerGroup
	orderConsumer platform_kafka.Consumer
	orderDecoder  kafka_decoder.OrderPaidDecoder

	syncProducer  sarama.SyncProducer
	orderProducer platform_kafka.Producer
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if c.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(config.GetConfig().Kafka.Brokers(), config.GetConfig().OrderAssembledConsumer.GroupID(), config.GetConfig().OrderAssembledConsumer.Config())
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("consumer_group", func(_ context.Context) error {
			return consumerGroup.Close()
		})
		c.consumerGroup = consumerGroup
	}
	return c.consumerGroup
}

func (c *diContainer) OrderConsumer() platform_kafka.Consumer {
	if c.orderConsumer == nil {
		c.orderConsumer = platform_kafka_consumer.NewConsumer(c.ConsumerGroup(), []string{config.GetConfig().OrderAssembledConsumer.Topic()}, logger.Logger())
	}
	return c.orderConsumer
}

func (c *diContainer) ConsumerService() service.ConsumerService {
	if c.consumerService == nil {
		c.consumerService = order_consumer.NewService(c.OrderConsumer(), c.OrderDecoder(), c.ProducerService())
	}
	return c.consumerService
}

func (c *diContainer) OrderDecoder() kafka_decoder.OrderPaidDecoder {
	if c.orderDecoder == nil {
		c.orderDecoder = decoder.NewOrderPaidDecoder()
	}
	return c.orderDecoder
}

func (c *diContainer) ProducerService() service.ProducerService {
	if c.producerService == nil {
		c.producerService = order_producer.NewService(c.OrderProducer())
	}
	return c.producerService
}

func (c *diContainer) OrderProducer() platform_kafka.Producer {
	if c.orderProducer == nil {
		p, err := platform_kafka_producer.NewProducer(c.SyncProducer(), config.GetConfig().OrderAssembledProducer.Topic(), logger.Logger())
		if err != nil {
			panic(fmt.Sprintf("failed to create order producer: %s\n", err.Error()))
		}
		c.orderProducer = p
	}
	return c.orderProducer
}

func (c *diContainer) SyncProducer() sarama.SyncProducer {
	if c.syncProducer == nil {
		p, err := sarama.NewSyncProducer(config.GetConfig().Kafka.Brokers(), config.GetConfig().OrderAssembledProducer.Config())
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("sync_producer", func(_ context.Context) error {
			return p.Close()
		})

		c.syncProducer = p
	}
	return c.syncProducer
}
