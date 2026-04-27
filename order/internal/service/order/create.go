package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (s *service) Create(ctx context.Context, info model.CreateOrderInfo) (model.Order, error) {
	partUuids := info.PartUUIDs()

	parts, err := s.inventoryClient.ListParts(ctx, partUuids)
	if err != nil {
		return model.Order{}, err
	}

	if len(parts) != len(partUuids) {
		return model.Order{}, errs.ErrPartNotFound
	}

	var totalPrice int64

	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, errs.ErrOutOfStock
		}

		totalPrice += part.Price
	}

	order := model.Order{
		OrderUUID:  uuid.New(),
		HullUUID:   info.HullUUID,
		EngineUUID: info.EngineUUID,
		ShieldUUID: info.ShieldUUID,
		WeaponUUID: info.WeaponUUID,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	err = s.repository.Create(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
