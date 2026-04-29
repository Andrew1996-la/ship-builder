package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/client/grpc/payment/v1/converter"
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

type Client struct {
	client paymentv1.PaymentServiceClient
}

func New(client paymentv1.PaymentServiceClient) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) PayOrder(
	ctx context.Context,
	orderUUID uuid.UUID,
	paymentMethod model.PaymentMethod,
) (uuid.UUID, error) {
	resp, err := c.client.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: converter.ToProtoPaymentMethod(paymentMethod),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ через сервис оплаты: %w", err)
	}

	transactionUUID, err := converter.ToModelTransactionUUID(resp)
	if err != nil {
		return uuid.Nil, fmt.Errorf("преобразовать ответ сервиса оплаты: %w", err)
	}

	return transactionUUID, nil
}
