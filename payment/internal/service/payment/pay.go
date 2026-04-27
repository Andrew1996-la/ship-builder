package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/payment/internal/errors"
	"github.com/Andrew1996-la/ship-builder/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, info model.PayRequest) (model.Payment, error) {
	if info.PaymentMethod == model.PaymentMethodUnspecified {
		return model.Payment{}, errs.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.New()

	slog.Info(
		"оплата выполнена",
		"order_uuid", info.OrderUUID.String(),
		"transaction_uuid", transactionUUID.String(),
	)

	return model.Payment{
		OrderUUID:       info.OrderUUID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   info.PaymentMethod,
	}, nil
}
