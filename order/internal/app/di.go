package app

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/H1dEx/ms-rocket/order/internal/api/order/v1"
	grpcClient "github.com/H1dEx/ms-rocket/order/internal/client/grpc"
	inventoryCli "github.com/H1dEx/ms-rocket/order/internal/client/grpc/inventory/v1"
	paymentCli "github.com/H1dEx/ms-rocket/order/internal/client/grpc/payment/v1"
	"github.com/H1dEx/ms-rocket/order/internal/config"
	kafka_decoder "github.com/H1dEx/ms-rocket/order/internal/converter/kafka"
	"github.com/H1dEx/ms-rocket/order/internal/converter/kafka/decoder"
	"github.com/H1dEx/ms-rocket/order/internal/repository"
	orderRepo "github.com/H1dEx/ms-rocket/order/internal/repository/order"
	"github.com/H1dEx/ms-rocket/order/internal/service"
	consumerService "github.com/H1dEx/ms-rocket/order/internal/service/consumer/order_consumer"
	orderService "github.com/H1dEx/ms-rocket/order/internal/service/order"
	producerService "github.com/H1dEx/ms-rocket/order/internal/service/producer/order_producer"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	k_consumer "github.com/H1dEx/ms-rocket/platform/pkg/kafka/consumer"
	k_producer "github.com/H1dEx/ms-rocket/platform/pkg/kafka/producer"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	"github.com/H1dEx/ms-rocket/platform/pkg/migrator"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API   orderV1.Handler
	orderService service.OrderService
	orderRepo    repository.OrderRepository

	inventoryClient grpcClient.InventoryClient
	inventoryConn   *grpc.ClientConn
	paymentClient   grpcClient.PaymentClient
	paymentConn     *grpc.ClientConn

	orderProducer service.ProducerService
	orderConsumer service.ConsumerService

	kafkaProducer kafka.Producer
	kafkaConsumer kafka.Consumer
	consumerGroup sarama.ConsumerGroup
	orderDecoder  kafka_decoder.OrderPaidDecoder

	syncProducer sarama.SyncProducer

	postgresConn *pgxpool.Pool
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if c.orderV1API == nil {
		c.orderV1API = orderApi.NewOrderAPI(c.OrderService(ctx))
	}
	return c.orderV1API
}

func (c *diContainer) OrderService(ctx context.Context) service.OrderService {
	if c.orderService == nil {
		c.orderService = orderService.NewOrderService(c.OrderRepo(ctx), c.InventoryClient(ctx), c.PaymentClient(ctx), c.OrderProducer(ctx))
	}
	return c.orderService
}

func (c *diContainer) OrderRepo(ctx context.Context) repository.OrderRepository {
	if c.orderRepo == nil {
		c.orderRepo = orderRepo.NewOrderRepository(c.PostgresConn(ctx))
	}
	return c.orderRepo
}

func (c *diContainer) PostgresConn(ctx context.Context) *pgxpool.Pool {
	if c.postgresConn == nil {
		conn, err := pgxpool.New(ctx, config.GetConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to database: %s", err.Error()))
		}
		closer.AddNamed("PostgreSQL connection", func(_ context.Context) error {
			conn.Close()
			return nil
		})

		err = conn.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to ping database: %s", err.Error()))
		}
		sqlDB := stdlib.OpenDB(*conn.Config().ConnConfig)

		closer.AddNamed("PostgreSQL converted connection", func(ctx context.Context) error {
			if cerr := sqlDB.Close(); cerr != nil {
				logger.Error(ctx, "failed to close database", zap.Error(cerr))
				return cerr
			}
			return nil
		})

		orderDir, err := findOrderDir()
		if err != nil {
			panic(fmt.Errorf("failed to find order directory: %s", err.Error()))
		}

		migrationDir := filepath.Join(orderDir, config.GetConfig().Postgres.MigrationDir())
		migrator := migrator.NewMigrator(sqlDB, migrationDir)
		err = migrator.Up()
		if err != nil {
			panic(fmt.Errorf("failed to migrate database: %s", err.Error()))
		}
		c.postgresConn = conn
	}
	return c.postgresConn
}

func (c *diContainer) InventoryClient(ctx context.Context) grpcClient.InventoryClient {
	if c.inventoryClient == nil {
		cli := inventoryV1.NewInventoryServiceClient(c.InventoryConn(ctx))
		c.inventoryClient = inventoryCli.NewInventoryClient(cli)
	}
	return c.inventoryClient
}

func (c *diContainer) InventoryConn(_ context.Context) *grpc.ClientConn {
	if c.inventoryConn == nil {
		inventoryConn, err := grpc.NewClient(
			config.GetConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to Inventory gRPC: %s", err.Error()))
		}
		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			if cerr := inventoryConn.Close(); cerr != nil {
				logger.Error(ctx, "failed to close connect", zap.Error(cerr))
				return cerr
			}
			return nil
		})
		c.inventoryConn = inventoryConn
	}
	return c.inventoryConn
}

func (c *diContainer) PaymentClient(ctx context.Context) grpcClient.PaymentClient {
	if c.paymentClient == nil {
		cli := paymentV1.NewPaymentServiceClient(c.PaymentConn(ctx))
		c.paymentClient = paymentCli.NewPaymentClient(cli)
	}
	return c.paymentClient
}

func (c *diContainer) PaymentConn(_ context.Context) *grpc.ClientConn {
	if c.paymentConn == nil {
		paymentConn, err := grpc.NewClient(
			config.GetConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to Payment gRPC: %s", err.Error()))
		}
		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			if cerr := paymentConn.Close(); cerr != nil {
				logger.Error(ctx, "failed to close connect", zap.Error(cerr))
				return cerr
			}
			return nil
		})
		c.paymentConn = paymentConn
	}
	return c.paymentConn
}

func findOrderDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for filepath.Base(dir) != "order" {
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("order not found")
		}
		dir = parent
	}
	return dir, nil
}

func (c *diContainer) OrderProducer(ctx context.Context) service.ProducerService {
	if c.orderProducer == nil {
		c.orderProducer = producerService.NewService(c.KafkaProducer(ctx))
	}
	return c.orderProducer
}

func (c *diContainer) KafkaProducer(ctx context.Context) kafka.Producer {
	if c.kafkaProducer == nil {
		p, err := k_producer.NewProducer(c.SyncProducer(ctx), config.GetConfig().OrderPaidProducer.Topic(), logger.Logger())
		if err != nil {
			panic(fmt.Errorf("failed to create kafka producer: %s", err.Error()))
		}
		c.kafkaProducer = p
	}
	return c.kafkaProducer
}

func (c *diContainer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if c.syncProducer == nil {
		p, err := sarama.NewSyncProducer(config.GetConfig().Kafka.Brokers(), config.GetConfig().OrderPaidProducer.Config())
		if err != nil {
			panic(fmt.Errorf("failed to create sync producer: %s", err.Error()))
		}
		closer.AddNamed("sync_producer", func(ctx context.Context) error {
			return p.Close()
		})
		c.syncProducer = p
	}
	return c.syncProducer
}

func (c *diContainer) OrderConsumer(ctx context.Context) service.ConsumerService {
	if c.orderConsumer == nil {
		c.orderConsumer = consumerService.NewService(c.OrderRepo(ctx), c.KafkaConsumer(ctx), c.OrderDecoder(ctx))
	}
	return c.orderConsumer
}

func (c *diContainer) OrderDecoder(ctx context.Context) kafka_decoder.OrderPaidDecoder {
	if c.orderDecoder == nil {
		c.orderDecoder = decoder.NewOrderPaidDecoder()
	}
	return c.orderDecoder
}

func (c *diContainer) KafkaConsumer(ctx context.Context) kafka.Consumer {
	if c.kafkaConsumer == nil {
		c.kafkaConsumer = k_consumer.NewConsumer(c.ConsumerGroup(ctx), []string{config.GetConfig().OrderAssembledConsumer.Topic()}, logger.Logger())
	}
	return c.kafkaConsumer
}

func (c *diContainer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if c.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(config.GetConfig().Kafka.Brokers(), config.GetConfig().OrderAssembledConsumer.GroupID(), config.GetConfig().OrderAssembledConsumer.Config())
		if err != nil {
			panic(fmt.Errorf("failed to create consumer group: %s", err.Error()))
		}
		closer.AddNamed("consumer_group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})
		c.consumerGroup = consumerGroup
	}
	return c.consumerGroup
}
