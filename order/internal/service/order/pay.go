package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (s *service) Pay(ctx context.Context, info model.PayOrderInfo) (model.Order, error) {
	order, err := s.repository.Get(ctx, info.OrderUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ для оплаты: %w", err)
	}

	if order.Status == model.OrderStatusPaid {
		return model.Order{}, errs.ErrOrderAlreadyPaid
	}

	if order.Status == model.OrderStatusCancelled {
		return model.Order{}, errs.ErrOrderCancelled
	}

	transactionUUIDRaw, err := s.paymentClient.PayOrder(ctx, info.OrderUUID, info.PaymentMethod)
	if err != nil {
		return model.Order{}, fmt.Errorf("оплатить заказ: %w", err)
	}

	transactionUUID, err := uuid.Parse(transactionUUIDRaw)
	if err != nil {
		return model.Order{}, fmt.Errorf("разобрать UUID транзакции оплаты: %w", err)
	}

	if transactionUUID == uuid.Nil {
		return model.Order{}, fmt.Errorf("пустой UUID транзакции оплаты")
	}

	order.Status = model.OrderStatusPaid
	order.PaymentMethod = &info.PaymentMethod
	order.TransactionUUID = &transactionUUID

	err = s.repository.Update(ctx, order)
	if err != nil {
		return model.Order{}, fmt.Errorf("сохранить оплаченный заказ: %w", err)
	}

	return order, nil
}
