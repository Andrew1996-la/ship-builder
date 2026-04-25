package part

import (
	"context"

	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
)

func (s service) Get(ctx context.Context, uuid string) (model.Part, error) {
	return s.repository.Get(ctx, uuid)
}
