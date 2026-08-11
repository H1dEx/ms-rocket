package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/H1dEx/ms-rocket/notification/internal/converter"
	"github.com/H1dEx/ms-rocket/notification/internal/converter/kafka"
	"github.com/H1dEx/ms-rocket/notification/internal/model"
	events_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/events/v1"
)

var _ kafka.OrderPaidDecoder = (*orderPaidDecoder)(nil)

type orderPaidDecoder struct{}

func NewOrderPaidDecoder() *orderPaidDecoder {
	return &orderPaidDecoder{}
}

func (d *orderPaidDecoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb events_v1.OrderPaid
	err := proto.Unmarshal(data, &pb)
	if err != nil {
		return model.OrderPaidEvent{}, err
	}
	event := model.OrderPaidEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   converter.PaymentMethodToModel(pb.PaymentMethod),
		TransactionUUID: pb.TransactionUuid,
	}
	return event, nil
}
