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

func TestCancelOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()

	tests := []struct {
		name         string
		setupMock    func(service *mocks.OrderService)
		expectedType any
		expectedErr  error
	}{
		{
			name: "успешный сценарий",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Cancel(ctx, orderUUID).
					Return(model.Order{OrderUUID: orderUUID, Status: model.OrderStatusCancelled}, nil)
			},
			expectedType: &orderv1.CancelOrderResponse{},
		},
		{
			name: "не найдено",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Cancel(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expectedType: &orderv1.CancelOrderNotFound{},
		},
		{
			name: "уже оплачено",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Cancel(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderAlreadyPaid)
			},
			expectedType: &orderv1.CancelOrderConflict{},
		},
		{
			name: "отменено",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Cancel(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderCancelled)
			},
			expectedType: &orderv1.CancelOrderConflict{},
		},
		{
			name: "ошибка сервиса",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Cancel(ctx, orderUUID).
					Return(model.Order{}, assert.AnError)
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewOrderService(t)
			tt.setupMock(service)

			api := orderapi.New(service)

			resp, err := api.CancelOrder(ctx, orderv1.CancelOrderParams{OrderUUID: orderUUID})

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.IsType(t, tt.expectedType, resp)
		})
	}
}
