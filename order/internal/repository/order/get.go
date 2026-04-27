package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	"github.com/Andrew1996-la/ship-builder/order/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, orderUuid uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUuid]
	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return converter.RepositoryToModel(order), nil
}
