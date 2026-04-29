package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	orderservice "github.com/Andrew1996-la/ship-builder/order/internal/service/order"
	"github.com/Andrew1996-la/ship-builder/order/internal/service/order/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	hullUUID := uuid.New()
	engineUUID := uuid.New()
	shieldUUID := uuid.New()
	weaponUUID := uuid.New()

	info := model.CreateOrderInfo{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		ShieldUUID: &shieldUUID,
		WeaponUUID: &weaponUUID,
	}

	partUUIDs := []uuid.UUID{hullUUID, engineUUID, shieldUUID, weaponUUID}
	parts := []model.Part{
		{UUID: hullUUID, Price: 100, StockQuantity: 1},
		{UUID: engineUUID, Price: 200, StockQuantity: 2},
		{UUID: shieldUUID, Price: 300, StockQuantity: 3},
		{UUID: weaponUUID, Price: 400, StockQuantity: 4},
	}

	tests := []struct {
		name        string
		setupMock   func(repository *mocks.Repository, inventoryClient *mocks.InventoryClient)
		expectedErr error
	}{
		{
			name: "успешный сценарий",
			setupMock: func(repository *mocks.Repository, inventoryClient *mocks.InventoryClient) {
				inventoryClient.EXPECT().
					ListParts(ctx, partUUIDs).
					Return(parts, nil)

				repository.EXPECT().
					Create(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.OrderUUID != uuid.Nil &&
							order.HullUUID == hullUUID &&
							order.EngineUUID == engineUUID &&
							order.ShieldUUID != nil && *order.ShieldUUID == shieldUUID &&
							order.WeaponUUID != nil && *order.WeaponUUID == weaponUUID &&
							order.TotalPrice == 1000 &&
							order.TransactionUUID == nil &&
							order.PaymentMethod == nil &&
							order.Status == model.OrderStatusPendingPayment &&
							!order.CreatedAt.IsZero()
					})).
					Return(nil)
			},
		},
		{
			name: "ошибка клиента склада",
			setupMock: func(repository *mocks.Repository, inventoryClient *mocks.InventoryClient) {
				inventoryClient.EXPECT().
					ListParts(ctx, partUUIDs).
					Return(nil, errs.ErrPartNotFound)
			},
			expectedErr: errs.ErrPartNotFound,
		},
		{
			name: "деталь не найдена",
			setupMock: func(repository *mocks.Repository, inventoryClient *mocks.InventoryClient) {
				inventoryClient.EXPECT().
					ListParts(ctx, partUUIDs).
					Return(parts[:3], nil)
			},
			expectedErr: errs.ErrPartNotFound,
		},
		{
			name: "нет на складе",
			setupMock: func(repository *mocks.Repository, inventoryClient *mocks.InventoryClient) {
				outOfStockParts := append([]model.Part(nil), parts...)
				outOfStockParts[2].StockQuantity = 0

				inventoryClient.EXPECT().
					ListParts(ctx, partUUIDs).
					Return(outOfStockParts, nil)
			},
			expectedErr: errs.ErrOutOfStock,
		},
		{
			name: "ошибка репозитория",
			setupMock: func(repository *mocks.Repository, inventoryClient *mocks.InventoryClient) {
				inventoryClient.EXPECT().
					ListParts(ctx, partUUIDs).
					Return(parts, nil)

				repository.EXPECT().
					Create(ctx, mock.AnythingOfType("model.Order")).
					Return(errs.ErrOrderNotFound)
			},
			expectedErr: errs.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := mocks.NewRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			tt.setupMock(repository, inventoryClient)

			service := orderservice.New(repository, inventoryClient, nil)

			actual, err := service.Create(ctx, info)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Equal(t, model.Order{}, actual)

				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, actual.OrderUUID)
			assert.Equal(t, hullUUID, actual.HullUUID)
			assert.Equal(t, engineUUID, actual.EngineUUID)
			require.NotNil(t, actual.ShieldUUID)
			assert.Equal(t, shieldUUID, *actual.ShieldUUID)
			require.NotNil(t, actual.WeaponUUID)
			assert.Equal(t, weaponUUID, *actual.WeaponUUID)
			assert.Equal(t, int64(1000), actual.TotalPrice)
			assert.Nil(t, actual.TransactionUUID)
			assert.Nil(t, actual.PaymentMethod)
			assert.Equal(t, model.OrderStatusPendingPayment, actual.Status)
			assert.False(t, actual.CreatedAt.IsZero())
		})
	}
}
