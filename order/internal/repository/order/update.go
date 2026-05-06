package order

import (
	"context"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	query := `
		UPDATE orders
		SET
			total_price = $2,
			status = $3,
			transaction_uuid = $4,
			payment_method = $5,
			updated_at = NOW()
		WHERE uuid = $1
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		order.OrderUUID,
		order.TotalPrice,
		order.Status,
		order.TransactionUUID,
		order.PaymentMethod,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}
