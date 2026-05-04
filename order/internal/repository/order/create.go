package order

import (
	"context"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	query := `
		INSERT INTO orders (
			uuid,
			total_price,
			status,
			transaction_uuid,
			payment_method,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(
		ctx,
		query,
		order.OrderUUID,
		order.TotalPrice,
		order.Status,
		order.TransactionUUID,
		order.PaymentMethod,
		order.CreatedAt,
	)
	if err != nil {
		return err
	}

	itemsQuery := `
	INSERT INTO order_items (
		order_uuid,
		part_uuid,
		part_type,
		price
	)
	VALUES ($1, $2, $3, $4)
`

	_, err = r.pool.Exec(ctx, itemsQuery, order.OrderUUID, order.HullUUID, "HULL", 0)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, itemsQuery, order.OrderUUID, order.EngineUUID, "ENGINE", 0)
	if err != nil {
		return err
	}

	if order.ShieldUUID != nil {
		_, err = r.pool.Exec(ctx, itemsQuery, order.OrderUUID, *order.ShieldUUID, "SHIELD", 0)
		if err != nil {
			return err
		}
	}

	if order.WeaponUUID != nil {
		_, err = r.pool.Exec(ctx, itemsQuery, order.OrderUUID, *order.WeaponUUID, "WEAPON", 0)
		if err != nil {
			return err
		}
	}

	return nil
}
