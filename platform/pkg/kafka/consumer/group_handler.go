package consumer

import (
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Middleware func(handler kafka.MessageHandler) kafka.MessageHandler

type groupHandler struct {
	handler kafka.MessageHandler
	logger  Logger
}

func NewGroupHandler(handler kafka.MessageHandler, logger Logger, middlewares ...Middleware) *groupHandler {
	for _, mw := range middlewares {
		handler = mw(handler)
	}
	return &groupHandler{
		handler: handler,
		logger:  logger,
	}
}

func (h *groupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *groupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				h.logger.Info(session.Context(), "messages channel closed")
				return nil
			}
			msg := kafka.Message{
				Key:            message.Key,
				Value:          message.Value,
				Topic:          message.Topic,
				Partition:      message.Partition,
				Offset:         message.Offset,
				Timestamp:      message.Timestamp,
				BlockTimestamp: message.BlockTimestamp,
			}
			if err := h.handler(session.Context(), msg); err != nil {
				h.logger.Error(session.Context(), "failed to handle message", zap.Error(err))
				continue
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			h.logger.Info(session.Context(), "consumer group session done")
			return nil
		}
	}
}
