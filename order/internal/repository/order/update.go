package order

import (
	"context"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	"github.com/Andrew1996-la/ship-builder/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	r.mu.RLock()
	_, ok := r.orders[order.OrderUUID]
	r.mu.RUnlock()

	if !ok {
		return errs.ErrOrderNotFound
	}

	repoOrder := converter.ToRepoOrder(order)

	r.mu.Lock()
	r.orders[order.OrderUUID] = repoOrder
	r.mu.Unlock()

	return nil
}
