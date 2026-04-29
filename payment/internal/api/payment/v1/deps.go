package v1

import (
	"context"

	"github.com/Andrew1996-la/ship-builder/payment/internal/model"
)

type PaymentService interface {
	// Pay выполнить оплату заказа.
	//
	// Возвращаемые ошибки:
	// - errs.ErrInvalidOrderUUID — если UUID заказа некорректный
	// - errs.ErrInvalidPaymentMethod — если метод оплаты некорректный
	Pay(ctx context.Context, req model.PayRequest) (model.Payment, error)
}
