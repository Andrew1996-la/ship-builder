package v1

import (
	"context"
	"errors"
	"net/http"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	info := model.PayOrderInfo{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: model.PaymentMethod(req.PaymentMethod),
	}

	order, err := a.service.Pay(ctx, info)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.PayOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "заказ не найден",
			}, nil
		case errors.Is(err, errs.ErrOrderAlreadyPaid):
			return &orderv1.PayOrderConflict{
				Code:    http.StatusConflict,
				Message: "заказ уже оплачен",
			}, nil
		case errors.Is(err, errs.ErrOrderCancelled):
			return &orderv1.PayOrderConflict{
				Code:    http.StatusConflict,
				Message: "заказ отменён",
			}, nil
		default:
			return &orderv1.PayOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "ошибка оплаты",
			}, nil
		}
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: *order.TransactionUUID,
	}, nil
}
