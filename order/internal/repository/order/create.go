package order

import (
	"context"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

// func (r *repository) Create(ctx context.Context, order model.Order) error {
// 	return r.txManager.Do(ctx, func(ctx context.Context) error {
// 		return r.create(ctx, order)
// 	})
// }

func (r *repository) Create(ctx context.Context, order model.Order) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

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

	_, err := db.Exec(
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

	_, err = db.Exec(ctx, itemsQuery, order.OrderUUID, order.HullUUID, "HULL", 0)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, itemsQuery, order.OrderUUID, order.EngineUUID, "ENGINE", 0)
	if err != nil {
		return err
	}

	if order.ShieldUUID != nil {
		_, err = db.Exec(ctx, itemsQuery, order.OrderUUID, *order.ShieldUUID, "SHIELD", 0)
		if err != nil {
			return err
		}
	}

	if order.WeaponUUID != nil {
		_, err = db.Exec(ctx, itemsQuery, order.OrderUUID, *order.WeaponUUID, "WEAPON", 0)
		if err != nil {
			return err
		}
	}

	return nil
}
