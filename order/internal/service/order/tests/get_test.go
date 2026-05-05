package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	orderservice "github.com/Andrew1996-la/ship-builder/order/internal/service/order"
	"github.com/Andrew1996-la/ship-builder/order/internal/service/order/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()
	expectedOrder := model.Order{
		OrderUUID: orderUUID,
		Items: []model.OrderItem{
			{PartUUID: uuid.New(), PartType: model.PartTypeHull, Price: 300},
			{PartUUID: uuid.New(), PartType: model.PartTypeEngine, Price: 200},
		},
		TotalPrice: 500,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	tests := []struct {
		name        string
		setupMock   func(repository *mocks.Repository)
		expected    model.Order
		expectedErr error
	}{
		{
			name: "успешный сценарий",
			setupMock: func(repository *mocks.Repository) {
				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(expectedOrder, nil)
			},
			expected: expectedOrder,
		},
		{
			name: "ошибка репозитория",
			setupMock: func(repository *mocks.Repository) {
				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected:    model.Order{},
			expectedErr: errs.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := mocks.NewRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			tt.setupMock(repository)

			service := orderservice.New(repository, inventoryClient, nil)

			actual, err := service.Get(ctx, orderUUID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Equal(t, tt.expected, actual)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
