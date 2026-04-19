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
	if req.PaymentMethod == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "payment_method обязателен")
	}

	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат order_uuid: %s", req.GetOrderUuid())
	}

	transactionUUID := uuid.New()

	slog.Info(
		"оплата выполнена",
		"order_uuid", orderUUID.String(),
		"transaction_uuid", transactionUUID.String(),
	)

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
