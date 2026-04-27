package order

import (
	"context"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (s *service) Pay(ctx context.Context, info model.PayOrderInfo) (model.Order, error) {
	order, err := s.repository.Get(ctx, info.OrderUUID)
	if err != nil {
		return model.Order{}, nil
	}

	if order.Status == model.OrderStatusPaid {
		return model.Order{}, errs.ErrOrderAlreadyPaid
	}

	if order.Status == model.OrderStatusCancelled {
		return model.Order{}, errs.ErrOrderCancelled
	}

	transactionUuid, err := s.paymentClient.PayOrder(ctx, info.OrderUUID, info.PaymentMethod)
	if err != nil {
		return model.Order{}, nil
	}

	order.Status = model.OrderStatusPaid
	order.PaymentMethod = &info.PaymentMethod
	order.TransactionUUID = &transactionUuid

	err = s.repository.Update(ctx, order)
	if err != nil {
		return model.Order{}, nil
	}

	return order, nil
}
