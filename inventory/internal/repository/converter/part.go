package converter

import (
	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/repository/record"
)

func RepositoryToModel(part record.Part) (model.Part, error) {
	uuid, err := uuid.Parse(part.UUID)
	if err != nil {
		return model.Part{}, err
	}

	return model.Part{
		UUID:          uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      model.PartType(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}, nil
}

func ModelToRepository(part model.Part) (record.Part, error) {
	return record.Part{
		UUID:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      string(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}, nil
}
