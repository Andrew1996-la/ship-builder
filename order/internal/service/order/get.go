package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid uuid.UUID) (model.Order, error) {
	return s.repository.Get(ctx, uuid)
}
