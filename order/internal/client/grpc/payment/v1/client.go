package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Andrew1996-la/ship-builder/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
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
) (string, error) {
	resp, err := c.client.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: converter.ToProtoPaymentMethod(paymentMethod),
	})
	if err != nil {
		return "", mapPayOrderError(err)
	}

	return resp.GetTransactionUuid(), nil
}

func mapPayOrderError(err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return fmt.Errorf("оплатить заказ через сервис оплаты: %w", errors.Join(errs.ErrPaymentFailed, err))
	default:
		return fmt.Errorf("оплатить заказ через сервис оплаты: %w", err)
	}
}
