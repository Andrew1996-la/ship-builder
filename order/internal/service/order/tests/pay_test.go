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

type paymentClientFunc func(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (string, error)

func (f paymentClientFunc) PayOrder(
	ctx context.Context,
	orderUUID uuid.UUID,
	paymentMethod model.PaymentMethod,
) (string, error) {
	return f(ctx, orderUUID, paymentMethod)
}

func TestPay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()
	transactionUUID := uuid.New()
	paymentMethod := model.PaymentMethodCard

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
		setupMock   func(repository *mocks.Repository) paymentClientFunc
		expectedErr error
	}{
		{
			name: "успешный сценарий",
			setupMock: func(repository *mocks.Repository) paymentClientFunc {
				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(pendingOrder, nil)

				repository.EXPECT().
					Update(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.OrderUUID == orderUUID &&
							order.TransactionUUID != nil && *order.TransactionUUID == transactionUUID &&
							order.PaymentMethod != nil && *order.PaymentMethod == paymentMethod &&
							order.Status == model.OrderStatusPaid
					})).
					Return(nil)

				return func(
					ctx context.Context,
					actualOrderUUID uuid.UUID,
					actualPaymentMethod model.PaymentMethod,
				) (string, error) {
					assert.Equal(t, orderUUID, actualOrderUUID)
					assert.Equal(t, paymentMethod, actualPaymentMethod)

					return transactionUUID.String(), nil
				}
			},
		},
		{
			name: "заказ не найден",
			setupMock: func(repository *mocks.Repository) paymentClientFunc {
				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)

				return nil
			},
			expectedErr: errs.ErrOrderNotFound,
		},
		{
			name: "заказ уже оплачен",
			setupMock: func(repository *mocks.Repository) paymentClientFunc {
				paidOrder := pendingOrder
				paidOrder.Status = model.OrderStatusPaid

				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(paidOrder, nil)

				return nil
			},
			expectedErr: errs.ErrOrderAlreadyPaid,
		},
		{
			name: "заказ отменён",
			setupMock: func(repository *mocks.Repository) paymentClientFunc {
				cancelledOrder := pendingOrder
				cancelledOrder.Status = model.OrderStatusCancelled

				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(cancelledOrder, nil)

				return nil
			},
			expectedErr: errs.ErrOrderCancelled,
		},
		{
			name: "ошибка клиента оплаты",
			setupMock: func(repository *mocks.Repository) paymentClientFunc {
				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(pendingOrder, nil)

				return func(
					ctx context.Context,
					actualOrderUUID uuid.UUID,
					actualPaymentMethod model.PaymentMethod,
				) (string, error) {
					assert.Equal(t, orderUUID, actualOrderUUID)
					assert.Equal(t, paymentMethod, actualPaymentMethod)

					return "", errs.ErrOrderCancelled
				}
			},
			expectedErr: errs.ErrOrderCancelled,
		},
		{
			name: "ошибка обновления в репозитории",
			setupMock: func(repository *mocks.Repository) paymentClientFunc {
				repository.EXPECT().
					Get(ctx, orderUUID).
					Return(pendingOrder, nil)

				repository.EXPECT().
					Update(ctx, mock.AnythingOfType("model.Order")).
					Return(errs.ErrOrderNotFound)

				return func(
					ctx context.Context,
					actualOrderUUID uuid.UUID,
					actualPaymentMethod model.PaymentMethod,
				) (string, error) {
					assert.Equal(t, orderUUID, actualOrderUUID)
					assert.Equal(t, paymentMethod, actualPaymentMethod)

					return transactionUUID.String(), nil
				}
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
			paymentClient := tt.setupMock(repository)

			service := orderservice.New(repository, inventoryClient, paymentClient)

			actual, err := service.Pay(ctx, model.PayOrderInfo{
				OrderUUID:     orderUUID,
				PaymentMethod: paymentMethod,
			})

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Equal(t, model.Order{}, actual)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, model.OrderStatusPaid, actual.Status)
			require.NotNil(t, actual.TransactionUUID)
			assert.Equal(t, transactionUUID, *actual.TransactionUUID)
			require.NotNil(t, actual.PaymentMethod)
			assert.Equal(t, paymentMethod, *actual.PaymentMethod)
		})
	}
}
