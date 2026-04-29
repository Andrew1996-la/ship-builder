package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Andrew1996-la/ship-builder/inventory/internal/converter"
	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	filter := model.PartFilter{
		UUIDs:    req.GetUuids(),
		PartType: converter.ToModelPartType(req.GetPartType()),
	}

	parts, err := a.partService.List(ctx, filter)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, "детали не найдены")
		}
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.ToProtoParts(parts),
	}, nil
}
