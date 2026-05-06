package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/Andrew1996-la/ship-builder/payment/internal/errors"
)

func ErrorInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		return nil, mapError(err)
	}

	return resp, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, errs.ErrInvalidOrderUUID):
		return status.Error(codes.InvalidArgument, "неверный формат order_uuid")

	case errors.Is(err, errs.ErrInvalidPaymentMethod):
		return status.Error(codes.InvalidArgument, "payment_method обязателен")

	default:
		return status.Error(codes.Internal, "внутренняя ошибка сервера")
	}
}
