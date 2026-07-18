package order

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

func (s *service) GetOrderByID(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.repo.GetOrder(ctx, orderUUID)

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}