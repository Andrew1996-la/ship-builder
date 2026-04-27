package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

type OrderService interface {
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Create(ctx context.Context, info model.CreateOrderInfo) (model.Order, error)
	Pay(ctx context.Context, info model.PayOrderInfo) (model.Order, error)
	Cancel(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
}
