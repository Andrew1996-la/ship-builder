package part

import (
	"context"

	"github.com/jackc/pgx/v5"

	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
)

func (r *repository) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	if len(filter.UUIDs) > 0 {
		return r.listByUUIDs(ctx, filter.UUIDs)
	}

	return r.listByType(ctx, filter.PartType)
}

func (r *repository) listByUUIDs(ctx context.Context, uuids []string) ([]model.Part, error) {
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
		WHERE uuid = ANY($1)
		ORDER BY name
	`

	rows, err := r.pool.Query(ctx, query, uuids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parts, err := scanParts(rows)
	if err != nil {
		return nil, err
	}

	if len(parts) != len(uuids) {
		return nil, errs.ErrPartNotFound
	}

	return parts, nil
}

func (r *repository) listByType(ctx context.Context, partType model.PartType) ([]model.Part, error) {
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
		WHERE $1 = 'UNSPECIFIED' OR part_type = $1
		ORDER BY name
	`

	rows, err := r.pool.Query(ctx, query, partType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanParts(rows)
}

func scanParts(rows pgx.Rows) ([]model.Part, error) {
	parts := make([]model.Part, 0)

	for rows.Next() {
		var part model.Part

		err := rows.Scan(
			&part.UUID,
			&part.Name,
			&part.Description,
			&part.Price,
			&part.PartType,
			&part.StockQuantity,
			&part.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		parts = append(parts, part)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parts, nil
}
