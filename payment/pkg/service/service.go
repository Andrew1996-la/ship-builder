package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

// PaymentServer реализует gRPC сервис оплаты.
type PaymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
}

// PayOrder обрабатывает оплату заказа.
func (s *PaymentServer) PayOrder(
	_ context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	// 1. Проверить, что order_uuid не пустой → INVALID_ARGUMENT
	if req.OrderUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid обязателен")
	}

	// 2. Проверить, что payment_method != UNSPECIFIED → INVALID_ARGUMENT
	if req.PaymentMethod == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "payment_method обязателен")
	}

	// 3. Проверить формат UUID → INVALID_ARGUMENT
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат order_uuid: %s", req.GetOrderUuid())
	}

	// 4. Сгенерировать transaction_uuid (UUID v4)
	transactionUUID := uuid.New()

	// 5. Вывести в лог: "оплата прошла успешно, order_uuid: X, transaction_uuid: Y"
	slog.Info(
		"оплата выполнена",
		"order_uuid", orderUUID.String(),
		"transaction_uuid", transactionUUID.String(),
	)

	// 6. Вернуть transaction_uuid
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
