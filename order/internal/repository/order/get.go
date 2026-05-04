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

	if err = r.loadItems(ctx, &order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *repository) loadItems(ctx context.Context, order *model.Order) error {
	query := `
		SELECT
			part_uuid,
			part_type
		FROM order_items
		WHERE order_uuid = $1
	`

	rows, err := r.pool.Query(ctx, query, order.OrderUUID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			partUUID uuid.UUID
			partType string
		)

		if scanErr := rows.Scan(&partUUID, &partType); scanErr != nil {
			return scanErr
		}

		switch partType {
		case "HULL":
			order.HullUUID = partUUID
		case "ENGINE":
			order.EngineUUID = partUUID
		case "SHIELD":
			order.ShieldUUID = &partUUID
		case "WEAPON":
			order.WeaponUUID = &partUUID
		}
	}

	return rows.Err()
}
