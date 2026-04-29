package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

type Repository interface {
	// Create сохранить новый заказ.
	//
	// Возвращаемые ошибки:
	// - error — при ошибке сохранения
	Create(ctx context.Context, order model.Order) error

	// Get получить заказ по UUID.
	//
	// Возвращаемые ошибки:
	// - errs.ErrOrderNotFound — если заказ не найден
	Get(ctx context.Context, uuid uuid.UUID) (model.Order, error)

	// Update обновить заказ.
	//
	// Возвращаемые ошибки:
	// - errs.ErrOrderNotFound — если заказ не найден
	// - error — при ошибке сохранения
	Update(ctx context.Context, order model.Order) error
}

// InventoryClient описывает клиент для сервиса inventory.
type InventoryClient interface {
	// ListParts получить список деталей по UUID.
	//
	// Возвращаемые ошибки:
	// - errs.ErrPartNotFound — если хотя бы одна деталь не найдена
	// - errs.ErrOutOfStock — если деталь отсутствует на складе
	// - error — при сетевых/внутренних ошибках
	ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error)
}

// PaymentClient описывает клиент для сервиса payment.
type PaymentClient interface {
	// PayOrder выполнить оплату заказа.
	//
	// Возвращаемые ошибки:
	// - errs.ErrPaymentFailed — если сервис оплаты отклонил запрос
	// - error — при сетевых/внутренних ошибках
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (string, error)
}
