package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

type OrderService interface {
	// Get получить заказ по UUID.
	//
	// Возвращаемые ошибки:
	// - errs.ErrOrderNotFound — если заказ не найден
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)

	// Create создать новый заказ.
	//
	// Возвращаемые ошибки:
	// - errs.ErrPartNotFound — если одна из деталей не найдена
	// - errs.ErrOutOfStock - если деталь отсутствует на складе
	Create(ctx context.Context, info model.CreateOrderInfo) (model.Order, error)

	// Pay оплатить заказ.
	//
	// Возвращаемые ошибки:
	// - errs.ErrOrderNotFound — если заказ не найден
	// - errs.ErrOrderAlreadyPaid — если заказ уже оплачен
	// - errs.ErrOrderCancelled — если заказ отменён
	// - errs.ErrInvalidPaymentMethod — если метод оплаты некорректный
	Pay(ctx context.Context, info model.PayOrderInfo) (model.Order, error)

	// Cancel отменить заказ.
	//
	// Возвращаемые ошибки:
	// - errs.ErrOrderNotFound — если заказ не найден
	// - errs.ErrOrderAlreadyPaid — если заказ уже оплачен
	// - errs.ErrOrderCancelled — если заказ уже отменён
	Cancel(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
}
