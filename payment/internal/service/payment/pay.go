package payment

import (
	"context"

	"github.com/google/uuid"
)

func (*service) PayOrder(context.Context) (string, error) {
	newUUID := uuid.NewString()

	return newUUID, nil
}
