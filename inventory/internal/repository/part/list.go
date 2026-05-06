package part

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
)

func (r *repository) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	builder := sq.
		Select(
			"uuid",
			"name",
			"description",
			"price",
			"part_type",
			"stock_quantity",
			"created_at",
		).
		From("parts").
		PlaceholderFormat(sq.Dollar)

	builder = applyFilters(builder, filter)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parts, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Part])
	if err != nil {
		return nil, err
	}

	if len(filter.UUIDs) > 0 {
		parts = orderByRequestedUUIDs(parts, filter.UUIDs)
	}

	return parts, nil
}

func orderByRequestedUUIDs(parts []model.Part, uuids []string) []model.Part {
	partsByUUID := make(map[string]model.Part, len(parts))
	for _, part := range parts {
		partsByUUID[part.UUID.String()] = part
	}

	orderedParts := make([]model.Part, 0, len(uuids))
	for _, rawUUID := range uuids {
		part, ok := partsByUUID[rawUUID]
		if !ok {
			continue
		}

		orderedParts = append(orderedParts, part)
	}

	return orderedParts
}

func applyFilters(
	builder sq.SelectBuilder,
	filter model.PartFilter,
) sq.SelectBuilder {
	if len(filter.UUIDs) > 0 {
		return builder.
			Where(sq.Eq{"uuid": filter.UUIDs}).
			OrderBy("name")
	}

	if filter.PartType != model.PartTypeUnspecified {
		return builder.
			Where(sq.Eq{"part_type": filter.PartType}).
			OrderBy("name")
	}

	return builder.OrderBy("name")
}
