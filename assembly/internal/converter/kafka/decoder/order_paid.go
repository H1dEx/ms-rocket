package decoder

import (
	"github.com/H1dEx/ms-rocket/assembly/internal/converter"
	"github.com/H1dEx/ms-rocket/assembly/internal/converter/kafka"
	"github.com/H1dEx/ms-rocket/assembly/internal/model"
	events_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

var _ kafka.OrderPaidDecoder = (*decoder)(nil)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb events_v1.OrderPaid
	err := proto.Unmarshal(data, &pb)
	if err != nil {
		return model.OrderPaidEvent{}, err
	}
	event := converter.OrderPaidEventToModel(&pb)

	return event, nil
}
