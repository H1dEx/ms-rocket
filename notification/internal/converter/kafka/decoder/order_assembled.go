package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/H1dEx/ms-rocket/notification/internal/converter/kafka"
	"github.com/H1dEx/ms-rocket/notification/internal/model"
	events_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/events/v1"
)

var _ kafka.OrderAssembledDecoder = (*orderAssembledDecoder)(nil)

type orderAssembledDecoder struct{}

func NewOrderAssembledDecoder() *orderAssembledDecoder {
	return &orderAssembledDecoder{}
}

func (d *orderAssembledDecoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb events_v1.ShipAssembled
	err := proto.Unmarshal(data, &pb)
	if err != nil {
		return model.ShipAssembledEvent{}, err
	}
	event := model.ShipAssembledEvent{
		EventUUID:        pb.EventUuid,
		OrderUUID:        pb.OrderUuid,
		UserUUID:         pb.UserUuid,
		BuildTimeSeconds: pb.BuildTimeSeconds,
	}
	return event, nil
}
