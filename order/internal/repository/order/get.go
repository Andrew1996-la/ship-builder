package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (r *repository) Get(ctx context.Context, orderUuid uuid.UUID) (model.Order, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	query := `
		SELECT
			o.uuid,
			o.total_price,
			o.status,
			o.transaction_uuid,
			o.payment_method,
			o.created_at,
			oi.uuid,
			oi.order_uuid,
			oi.part_uuid,
			oi.part_type,
			oi.price,
			oi.created_at
		FROM orders o
		JOIN order_items oi ON oi.order_uuid = o.uuid
		WHERE o.uuid = $1
		ORDER BY oi.created_at, oi.uuid
	`

	var order model.Order

	rows, err := db.Query(ctx, query, orderUuid)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()

	var found bool

	for rows.Next() {
		var (
			item     model.OrderItem
			partType string
		)

		err = rows.Scan(
			&order.OrderUUID,
			&order.TotalPrice,
			&order.Status,
			&order.TransactionUUID,
			&order.PaymentMethod,
			&order.CreatedAt,
			&item.UUID,
			&item.OrderUUID,
			&item.PartUUID,
			&partType,
			&item.Price,
			&item.CreatedAt,
		)
		if err != nil {
			return model.Order{}, err
		}

		item.PartType = model.PartType(partType)
		order.Items = append(order.Items, item)
		found = true
	}

	if err = rows.Err(); err != nil {
		return model.Order{}, err
	}

	if !found {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return order, nil
}
