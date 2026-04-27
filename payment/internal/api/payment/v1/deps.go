package v1

import (
	"context"

	"github.com/Andrew1996-la/ship-builder/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, req model.PayRequest) (model.Payment, error)
}
