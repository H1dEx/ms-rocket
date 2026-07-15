package service

import (
	"context"
)

type PaymentService interface {
	PayOrder(context.Context) (string, error)
}
