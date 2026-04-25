package v1

import (
	"context"

	"github.com/Andrew1996-la/ship-builder/inventory/internal/converter"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
)

func (a *api) List(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	filter := model.PartFilter{
		UUIDs:    req.GetUuids(),
		PartType: converter.ProtoPartTypeToModel(req.GetPartType()),
	}

	parts, err := a.partService.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	respParts := make([]*inventoryv1.Part, 0, len(parts))

	for _, p := range parts {
		respParts = append(respParts, converter.ModelToProtoPart(p))
	}

	return &inventoryv1.ListPartsResponse{
		Parts: respParts,
	}, nil
}
