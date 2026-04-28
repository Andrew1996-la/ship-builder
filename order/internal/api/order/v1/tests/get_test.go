package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderapi "github.com/Andrew1996-la/ship-builder/order/internal/api/order/v1"
	"github.com/Andrew1996-la/ship-builder/order/internal/api/order/v1/mocks"
	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
)

func TestGetOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()
	shieldUUID := uuid.New()
	weaponUUID := uuid.New()
	transactionUUID := uuid.New()
	paymentMethod := model.PaymentMethodCard
	createdAt := time.Now()

	expectedOrder := model.Order{
		OrderUUID:       orderUUID,
		HullUUID:        uuid.New(),
		EngineUUID:      uuid.New(),
		ShieldUUID:      &shieldUUID,
		WeaponUUID:      &weaponUUID,
		TotalPrice:      1000,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          model.OrderStatusPaid,
		CreatedAt:       createdAt,
	}

	tests := []struct {
		name         string
		setupMock    func(service *mocks.OrderService)
		expectedType any
		expectedErr  error
	}{
		{
			name: "success",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Get(ctx, orderUUID).
					Return(expectedOrder, nil)
			},
			expectedType: &orderv1.OrderDto{},
		},
		{
			name: "not found",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expectedType: &orderv1.GetOrderNotFound{},
		},
		{
			name: "service error",
			setupMock: func(service *mocks.OrderService) {
				service.EXPECT().
					Get(ctx, orderUUID).
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

			resp, err := api.GetOrder(ctx, orderv1.GetOrderParams{OrderUUID: orderUUID})

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.IsType(t, tt.expectedType, resp)

			if dto, ok := resp.(*orderv1.OrderDto); ok {
				assert.Equal(t, expectedOrder.OrderUUID, dto.OrderUUID)
				assert.Equal(t, expectedOrder.HullUUID, dto.HullUUID)
				assert.Equal(t, expectedOrder.EngineUUID, dto.EngineUUID)
				assert.Equal(t, expectedOrder.TotalPrice, dto.TotalPrice)
				assert.Equal(t, orderv1.OrderStatusPAID, dto.Status)
				assert.Equal(t, createdAt, dto.CreatedAt)

				actualShieldUUID, ok := dto.ShieldUUID.Get()
				require.True(t, ok)
				assert.Equal(t, shieldUUID, actualShieldUUID)

				actualWeaponUUID, ok := dto.WeaponUUID.Get()
				require.True(t, ok)
				assert.Equal(t, weaponUUID, actualWeaponUUID)

				actualTransactionUUID, ok := dto.TransactionUUID.Get()
				require.True(t, ok)
				assert.Equal(t, transactionUUID, actualTransactionUUID)

				actualPaymentMethod, ok := dto.PaymentMethod.Get()
				require.True(t, ok)
				assert.Equal(t, orderv1.PaymentMethodCARD, actualPaymentMethod)
			}
		})
	}
}
