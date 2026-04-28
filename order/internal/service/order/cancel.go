package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUuid uuid.UUID) (model.Order, error) {
	order, err := s.repository.Get(ctx, orderUuid)
	if err != nil {
		return model.Order{}, errs.ErrOrderNotFound
	}

	if order.Status == model.OrderStatusPaid {
		return model.Order{}, errs.ErrOrderAlreadyPaid
	}

	if order.Status == model.OrderStatusCancelled {
		return model.Order{}, errs.ErrOrderCancelled
	}

	order.Status = model.OrderStatusCancelled

	err = s.repository.Update(ctx, order)
	if err != nil {
		return model.Order{}, fmt.Errorf("сохранить отменённый заказ: %w", err)
	}

	return order, nil
}
