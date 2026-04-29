package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	for _, rawUUID := range filter.UUIDs {
		id, err := uuid.Parse(rawUUID)
		if err != nil {
			return nil, errs.ErrInvalidUUID
		}

		if id == uuid.Nil {
			return nil, errs.ErrInvalidUUID
		}
	}

	return s.repository.List(ctx, filter)
}
