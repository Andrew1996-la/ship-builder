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

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()
	hullUUID := uuid.New()
	engineUUID := uuid.New()
	shieldUUID := uuid.New()
	weaponUUID := uuid.New()

	req := &orderv1.CreateOrderRequest{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		ShieldUUID: orderv1.NewOptNilUUID(shieldUUID),
		WeaponUUID: orderv1.NewOptNilUUID(weaponUUID),
	}

	expectedInfo := model.CreateOrderInfo{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		ShieldUUID: &shieldUUID,
		WeaponUUID: &weaponUUID,
	}

	tests := []struct {
		name         string
		setupMock    func(service *mocks.OrderService)
		expectedType any
	}{
		{
			name: "success",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Create(ctx, expectedInfo).
					Return(model.Order{
						OrderUUID:  orderUUID,
						TotalPrice: 1000,
					}, nil)
			},
			expectedType: &orderv1.CreateOrderResponse{},
		},
		{
			name: "part not found",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Create(ctx, expectedInfo).
					Return(model.Order{}, errs.ErrPartNotFound)
			},
			expectedType: &orderv1.CreateOrderNotFound{},
		},
		{
			name: "out of stock",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Create(ctx, expectedInfo).
					Return(model.Order{}, errs.ErrOutOfStock)
			},
			expectedType: &orderv1.CreateOrderConflict{},
		},
		{
			name: "internal error",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Create(ctx, expectedInfo).
					Return(model.Order{}, assert.AnError)
			},
			expectedType: &orderv1.CreateOrderInternalServerError{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewOrderService(t)
			tt.setupMock(service)

			api := orderapi.New(service)

			resp, err := api.CreateOrder(ctx, req)

			require.NoError(t, err)
			require.IsType(t, tt.expectedType, resp)

			if success, ok := resp.(*orderv1.CreateOrderResponse); ok {
				assert.Equal(t, orderUUID, success.OrderUUID)
				assert.Equal(t, int64(1000), success.TotalPrice)
			}
		})
	}
}
