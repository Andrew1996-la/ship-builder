package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		name         string
		req          *paymentv1.PayOrderRequest
		setupMock    func(service *mocks.PaymentService)
		expectedErr  error
		expectedCode codes.Code
	}{
		{
			name: "успешный сценарий",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			setupMock: func(service *mocks.PaymentService) {
				service.EXPECT().
					Pay(ctx, model.PayRequest{
						OrderUUID:     orderUUID.String(),
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
			name: "неверный UUID заказа",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     "bad-uuid",
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			setupMock: func(service *mocks.PaymentService) {
				service.EXPECT().
					Pay(ctx, model.PayRequest{
						OrderUUID:     "bad-uuid",
						PaymentMethod: model.PaymentMethodCard,
					}).
					Return(model.Payment{}, errs.ErrInvalidOrderUUID)
			},
			expectedErr:  errs.ErrInvalidOrderUUID,
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "ошибка сервиса",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			},
			setupMock: func(service *mocks.PaymentService) {
				service.EXPECT().
					Pay(ctx, model.PayRequest{
						OrderUUID:     orderUUID.String(),
						PaymentMethod: model.PaymentMethodUnspecified,
					}).
					Return(model.Payment{}, errs.ErrInvalidPaymentMethod)
			},
			expectedErr:  errs.ErrInvalidPaymentMethod,
			expectedCode: codes.InvalidArgument,
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
				assert.Equal(t, tt.expectedCode, status.Code(err))
				assert.Equal(t, tt.expectedErr.Error(), status.Convert(err).Message())
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, transactionUUID.String(), resp.GetTransactionUuid())
		})
	}
}
