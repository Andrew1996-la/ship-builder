package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/payment/internal/errors"
	"github.com/Andrew1996-la/ship-builder/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, info model.PayRequest) (model.Payment, error) {
	orderUUID, err := uuid.Parse(info.OrderUUID)
	if err != nil {
		return model.Payment{}, errs.ErrInvalidOrderUUID
	}

	if orderUUID == uuid.Nil {
		return model.Payment{}, errs.ErrInvalidOrderUUID
	}

	if info.PaymentMethod == model.PaymentMethodUnspecified {
		return model.Payment{}, errs.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.New()

	slog.Info(
		"оплата выполнена",
		"order_uuid", orderUUID.String(),
		"transaction_uuid", transactionUUID.String(),
	)

	return model.Payment{
		OrderUUID:       orderUUID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   info.PaymentMethod,
	}, nil
}
