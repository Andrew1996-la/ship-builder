package order

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	if r.txManager == nil {
		return r.create(ctx, order)
	}

	return r.txManager.Do(ctx, func(ctx context.Context) error {
		return r.create(ctx, order)
	})
}

func (r *repository) create(ctx context.Context, order model.Order) error {
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
			uuid,
			order_uuid,
			part_uuid,
			part_type,
			price,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	batch := &pgx.Batch{}
	for _, item := range order.Items {
		batch.Queue(
			itemsQuery,
			item.UUID,
			item.OrderUUID,
			item.PartUUID,
			item.PartType.String(),
			item.Price,
			item.CreatedAt,
		)
	}

	results := db.SendBatch(ctx, batch)
	defer results.Close()

	for range order.Items {
		if _, err = results.Exec(); err != nil {
			return err
		}
	}

	return nil
}
