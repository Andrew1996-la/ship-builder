package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/Andrew1996-la/ship-builder/order/internal/converter"
	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.service.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return &orderv1.GetOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "заказ не найден",
			}, nil
		}
		return nil, err
	}

	return converter.ToDTO(order), nil
}
