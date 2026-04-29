package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderapi "github.com/Andrew1996-la/ship-builder/order/internal/api/order/v1"
	"github.com/Andrew1996-la/ship-builder/order/internal/api/order/v1/mocks"
	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()
	transactionUUID := uuid.New()
	req := &orderv1.PayOrderRequest{PaymentMethod: orderv1.PaymentMethodCARD}
	expectedInfo := model.PayOrderInfo{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCard,
	}

	tests := []struct {
		name         string
		setupMock    func(service *mocks.OrderService)
		expectedType any
	}{
		{
			name: "успешный сценарий",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Pay(ctx, expectedInfo).
					Return(model.Order{
						OrderUUID:       orderUUID,
						TransactionUUID: &transactionUUID,
					}, nil)
			},
			expectedType: &orderv1.PayOrderResponse{},
		},
		{
			name: "не найдено",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Pay(ctx, expectedInfo).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expectedType: &orderv1.PayOrderNotFound{},
		},
		{
			name: "уже оплачено",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Pay(ctx, expectedInfo).
					Return(model.Order{}, errs.ErrOrderAlreadyPaid)
			},
			expectedType: &orderv1.PayOrderConflict{},
		},
		{
			name: "отменено",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Pay(ctx, expectedInfo).
					Return(model.Order{}, errs.ErrOrderCancelled)
			},
			expectedType: &orderv1.PayOrderConflict{},
		},
		{
			name: "внутренняя ошибка",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Pay(ctx, expectedInfo).
					Return(model.Order{}, assert.AnError)
			},
			expectedType: &orderv1.PayOrderInternalServerError{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewOrderService(t)
			tt.setupMock(service)

			api := orderapi.New(service)

			resp, err := api.PayOrder(ctx, req, orderv1.PayOrderParams{OrderUUID: orderUUID})

			require.NoError(t, err)
			require.IsType(t, tt.expectedType, resp)

			if success, ok := resp.(*orderv1.PayOrderResponse); ok {
				assert.Equal(t, transactionUUID, success.TransactionUUID)
			}
		})
	}
}
