package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Andrew1996-la/ship-builder/payment/internal/converter"
	errs "github.com/Andrew1996-la/ship-builder/payment/internal/errors"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	info := converter.ToModelPayRequest(req)

	payment, err := a.service.Pay(ctx, info)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidOrderUUID):
			return nil, status.Error(codes.InvalidArgument, err.Error())

		case errors.Is(err, errs.ErrInvalidPaymentMethod):
			return nil, status.Error(codes.InvalidArgument, err.Error())

		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return converter.ToProtoPayOrderResponse(payment), nil
}
