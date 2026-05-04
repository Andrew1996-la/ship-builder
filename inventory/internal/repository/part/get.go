package part

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
)

func (r *repository) Get(ctx context.Context, id string) (model.Part, error) {
	query := `
		SELECT
			uuid,
			name,
			description,
			price,
			part_type,
			stock_quantity,
			created_at
		FROM parts
		WHERE uuid = $1
	`

	var part model.Part

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&part.UUID,
		&part.Name,
		&part.Description,
		&part.Price,
		&part.PartType,
		&part.StockQuantity,
		&part.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Part{}, errs.ErrPartNotFound
		}

		return model.Part{}, err
	}

	return part, nil
}
