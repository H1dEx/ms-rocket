package v1

import (
	"github.com/H1dEx/ms-rocket/order/internal/client/grpc"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

var _ grpc.PaymentClient = (*client)(nil)

type client struct {
	paymentClient paymentV1.PaymentServiceClient
}

func NewPaymentClient(paymentClient paymentV1.PaymentServiceClient) grpc.PaymentClient {
	return &client{
		paymentClient: paymentClient,
	}
}
