package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/Andrew1996-la/ship-builder/payment/internal/errors"
	"github.com/Andrew1996-la/ship-builder/payment/internal/model"
	paymentservice "github.com/Andrew1996-la/ship-builder/payment/internal/service/payment"
)

func TestPay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()

	tests := []struct {
		name        string
		info        model.PayRequest
		expectedErr error
	}{
		{
			name: "success card",
			info: model.PayRequest{
				OrderUUID:     orderUUID,
				PaymentMethod: model.PaymentMethodCard,
			},
		},
		{
			name: "success sbp",
			info: model.PayRequest{
				OrderUUID:     orderUUID,
				PaymentMethod: model.PaymentMethodSBP,
			},
		},
		{
			name: "invalid payment method",
			info: model.PayRequest{
				OrderUUID:     orderUUID,
				PaymentMethod: model.PaymentMethodUnspecified,
			},
			expectedErr: errs.ErrInvalidPaymentMethod,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := paymentservice.New()

			actual, err := service.Pay(ctx, tt.info)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Equal(t, model.Payment{}, actual)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.info.OrderUUID, actual.OrderUUID)
			assert.Equal(t, tt.info.PaymentMethod, actual.PaymentMethod)
			assert.NotEqual(t, uuid.Nil, actual.TransactionUUID)
		})
	}
}
