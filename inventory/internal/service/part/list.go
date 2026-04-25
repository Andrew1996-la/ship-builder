package part

import (
	"context"

	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	return s.repository.List(ctx, filter)
}
