package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	paymentapi "github.com/Andrew1996-la/ship-builder/payment/internal/api/payment/v1"
	"github.com/Andrew1996-la/ship-builder/payment/internal/api/payment/v1/mocks"
	errs "github.com/Andrew1996-la/ship-builder/payment/internal/errors"
	"github.com/Andrew1996-la/ship-builder/payment/internal/model"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()
	transactionUUID := uuid.New()

	tests := []struct {
		name        string
		req         *paymentv1.PayOrderRequest
		setupMock   func(service *mocks.PaymentService)
		expectedErr error
	}{
		{
			name: "success",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			setupMock: func(service *mocks.PaymentService) {
				service.EXPECT().
					Pay(ctx, model.PayRequest{
						OrderUUID:     orderUUID,
						PaymentMethod: model.PaymentMethodCard,
					}).
					Return(model.Payment{
						OrderUUID:       orderUUID,
						TransactionUUID: transactionUUID,
						PaymentMethod:   model.PaymentMethodCard,
					}, nil)
			},
		},
		{
			name: "invalid order uuid",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     "bad-uuid",
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			setupMock: func(service *mocks.PaymentService) {
				// service не должен вызываться, потому что converter упадёт раньше
			},
			expectedErr: errs.ErrInvalidOrderUUID,
		},
		{
			name: "service error",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			},
			setupMock: func(service *mocks.PaymentService) {
				service.EXPECT().
					Pay(ctx, model.PayRequest{
						OrderUUID:     orderUUID,
						PaymentMethod: model.PaymentMethodUnspecified,
					}).
					Return(model.Payment{}, errs.ErrInvalidPaymentMethod)
			},
			expectedErr: errs.ErrInvalidPaymentMethod,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewPaymentService(t)
			tt.setupMock(service)

			api := paymentapi.New(service)

			resp, err := api.PayOrder(ctx, tt.req)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, transactionUUID.String(), resp.GetTransactionUuid())
		})
	}
}
