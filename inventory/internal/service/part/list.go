package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	for i, rawUUID := range filter.UUIDs {
		id, err := uuid.Parse(rawUUID)
		if err != nil {
			return nil, errs.ErrInvalidUUID
		}

		if id == uuid.Nil {
			return nil, errs.ErrInvalidUUID
		}

		filter.UUIDs[i] = id.String()
	}

	parts, err := s.repository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(filter.UUIDs) > 0 && len(filter.UUIDs) != len(parts) {
		return nil, errs.ErrPartNotFound
	}

	return parts, nil
}
