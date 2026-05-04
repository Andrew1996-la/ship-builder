package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (r *repository) Get(ctx context.Context, orderUuid uuid.UUID) (model.Order, error) {
	query := `
		SELECT
			uuid,
			total_price,
			status,
			transaction_uuid,
			payment_method,
			created_at
		FROM orders
		WHERE uuid = $1
	`

	var order model.Order

	err := r.pool.QueryRow(ctx, query, orderUuid).Scan(
		&order.OrderUUID,
		&order.TotalPrice,
		&order.Status,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, errs.ErrOrderNotFound
		}

		return model.Order{}, err
	}

	return order, nil
}
