package v1

import (
	"context"

	"github.com/Andrew1996-la/ship-builder/payment/internal/converter"
	errs "github.com/Andrew1996-la/ship-builder/payment/internal/errors"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	info, err := converter.PayOrderRequestToModel(req)
	if err != nil {
		return nil, errs.ErrInvalidOrderUUID
	}

	payment, err := a.service.Pay(ctx, info)
	if err != nil {
		return nil, err
	}

	return converter.ModelToPayOrderResponse(payment), nil
}
