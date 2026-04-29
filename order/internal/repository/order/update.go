package order

import (
	"context"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	"github.com/Andrew1996-la/ship-builder/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.OrderUUID]; !ok {
		return errs.ErrOrderNotFound
	}

	r.orders[order.OrderUUID] = converter.ToRepoOrder(order)

	return nil
}
