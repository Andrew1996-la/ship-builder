package v1

import (
	"context"
	"errors"
	"net/http"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	info := model.CreateOrderInfo{
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
	}

	if shieldUuid, ok := req.ShieldUUID.Get(); ok {
		info.ShieldUUID = &shieldUuid
	}

	if weaponUuid, ok := req.WeaponUUID.Get(); ok {
		info.WeaponUUID = &weaponUuid
	}

	order, err := a.service.Create(ctx, info)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrPartNotFound):
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "деталь не найдена",
			}, nil
		case errors.Is(err, errs.ErrOutOfStock):
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: "деталь отсутсвует на складе",
			}, nil
		default:
			return &orderv1.CreateOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервиса",
			}, nil
		}
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
