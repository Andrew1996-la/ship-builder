package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
)

func ErrorInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}

	return resp, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, errs.ErrInvalidUUID):
		return status.Error(codes.InvalidArgument, "неверный формат UUID")

	case errors.Is(err, errs.ErrPartNotFound):
		return status.Error(codes.NotFound, "деталь не найдена")

	default:
		return status.Error(codes.Internal, "внутренняя ошибка сервиса")
	}
}
