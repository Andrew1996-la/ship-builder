package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	orderservice "github.com/Andrew1996-la/ship-builder/order/internal/service/order"
	"github.com/Andrew1996-la/ship-builder/order/internal/service/order/mocks"
)

func TestCancel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()

	pendingOrder := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
		TotalPrice: 500,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	tests := []struct {
		name        string
		setupMock   func(repository *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(repository *mocks.Repository) {
				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(pendingOrder, nil)

				repository.EXPECT().
					Update(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.OrderUUID == orderUUID &&
							order.Status == model.OrderStatusCancelled
					})).
					Return(nil)
			},
		},
		{
			name: "order not found",
			setupMock: func(repository *mocks.Repository) {
				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expectedErr: errs.ErrOrderNotFound,
		},
		{
			name: "order already paid",
			setupMock: func(repository *mocks.Repository) {
				paidOrder := pendingOrder
				paidOrder.Status = model.OrderStatusPaid

				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(paidOrder, nil)
			},
			expectedErr: errs.ErrOrderAlreadyPaid,
		},
		{
			name: "order cancelled",
			setupMock: func(repository *mocks.Repository) {
				cancelledOrder := pendingOrder
				cancelledOrder.Status = model.OrderStatusCancelled

				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(cancelledOrder, nil)
			},
			expectedErr: errs.ErrOrderCancelled,
		},
		{
			name: "repository update error",
			setupMock: func(repository *mocks.Repository) {
				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(pendingOrder, nil)

				repository.EXPECT().
					Update(ctx, mock.AnythingOfType("model.Order")).
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
			paymentClient := mocks.NewPaymentClient(t)
			tt.setupMock(repository)

			service := orderservice.New(repository, inventoryClient, paymentClient)

			actual, err := service.Cancel(ctx, orderUUID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Equal(t, model.Order{}, actual)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, model.OrderStatusCancelled, actual.Status)
		})
	}
}
