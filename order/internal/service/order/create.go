package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (s *service) Create(ctx context.Context, info model.CreateOrderInfo) (model.Order, error) {
	partUuids := info.PartUUIDs()

	parts, err := s.inventoryClient.ListParts(ctx, partUuids)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить детали для создания заказа: %w", err)
	}

	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, errs.ErrOutOfStock
		}
	}

	orderUUID := uuid.New()
	createdAt := time.Now()

	items, totalPrice, err := buildOrderItems(info, parts, orderUUID, createdAt)
	if err != nil {
		return model.Order{}, err
	}

	order := model.Order{
		OrderUUID:  orderUUID,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  createdAt,
	}

	if err = s.repository.Create(ctx, order); err != nil {
		return model.Order{}, fmt.Errorf("сохранить созданный заказ: %w", err)
	}

	return order, nil
}

func buildOrderItems(
	info model.CreateOrderInfo,
	parts []model.Part,
	orderUUID uuid.UUID,
	createdAt time.Time,
) ([]model.OrderItem, int64, error) {
	partsByUUID := make(map[uuid.UUID]model.Part, len(parts))
	for _, part := range parts {
		partsByUUID[part.UUID] = part
	}

	requestedItems := []struct {
		partUUID uuid.UUID
		partType model.PartType
	}{
		{partUUID: info.HullUUID, partType: model.PartTypeHull},
		{partUUID: info.EngineUUID, partType: model.PartTypeEngine},
	}

	if info.ShieldUUID != nil {
		requestedItems = append(requestedItems, struct {
			partUUID uuid.UUID
			partType model.PartType
		}{partUUID: *info.ShieldUUID, partType: model.PartTypeShield})
	}

	if info.WeaponUUID != nil {
		requestedItems = append(requestedItems, struct {
			partUUID uuid.UUID
			partType model.PartType
		}{partUUID: *info.WeaponUUID, partType: model.PartTypeWeapon})
	}

	items := make([]model.OrderItem, 0, len(requestedItems))
	var totalPrice int64

	for _, requestedItem := range requestedItems {
		part, ok := partsByUUID[requestedItem.partUUID]
		if !ok {
			return nil, 0, errs.ErrPartNotFound
		}

		items = append(items, model.OrderItem{
			UUID:      uuid.New(),
			OrderUUID: orderUUID,
			PartUUID:  requestedItem.partUUID,
			PartType:  requestedItem.partType,
			Price:     part.Price,
			CreatedAt: createdAt,
		})
		totalPrice += part.Price
	}

	return items, totalPrice, nil
}
