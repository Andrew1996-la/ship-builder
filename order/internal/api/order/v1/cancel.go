package v1

import (
	"context"
	"errors"
	"net/http"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	_, err := a.service.Cancel(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.CancelOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "заказ не найден",
			}, nil
		case errors.Is(err, errs.ErrOrderAlreadyPaid):
			return &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: "заказ уже оплачен",
			}, nil
		case errors.Is(err, errs.ErrOrderCancelled):
			return &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: "заказ уже отменён",
			}, nil
		default:
			return nil, err
		}
	}

	return &orderv1.CancelOrderResponse{}, nil
}
