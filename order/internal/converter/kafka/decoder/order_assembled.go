package decoder

import (
	"github.com/H1dEx/ms-rocket/order/internal/converter"
	"github.com/H1dEx/ms-rocket/order/internal/converter/kafka"
	"github.com/H1dEx/ms-rocket/order/internal/model"
	events_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

var _ kafka.OrderPaidDecoder = (*decoder)(nil)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb events_v1.ShipAssembled
	err := proto.Unmarshal(data, &pb)
	if err != nil {
		return model.ShipAssembledEvent{}, err
	}
	event := converter.ShipAssembledEventToModel(&pb)
	return event, nil
}