package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
)

func ToRawUUIDs(uuids []uuid.UUID) []string {
	rawUUIDs := make([]string, 0, len(uuids))

	for _, id := range uuids {
		rawUUIDs = append(rawUUIDs, id.String())
	}

	return rawUUIDs
}

func ToModelParts(parts []*inventoryv1.Part) ([]model.Part, error) {
	modelParts := make([]model.Part, 0, len(parts))

	for _, part := range parts {
		modelPart, err := ToModelPart(part)
		if err != nil {
			return nil, err
		}

		modelParts = append(modelParts, modelPart)
	}

	return modelParts, nil
}

func ToModelPart(part *inventoryv1.Part) (model.Part, error) {
	id, err := uuid.Parse(part.GetUuid())
	if err != nil {
		return model.Part{}, fmt.Errorf("разобрать UUID детали: %w", err)
	}

	return model.Part{
		UUID:          id,
		Name:          part.GetName(),
		PartType:      model.PartType(part.GetPartType().String()),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
	}, nil
}
