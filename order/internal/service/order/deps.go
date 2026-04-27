package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

type Repository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, uuid uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
