package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

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
		PaymentMethod: toProtoPaymentMethod(paymentMethod),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ через сервис оплаты: %w", err)
	}

	transactionUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, fmt.Errorf("разобрать UUID транзакции из ответа сервиса оплаты: %w", err)
	}

	return transactionUUID, nil
}

func toProtoPaymentMethod(paymentMethod model.PaymentMethod) paymentv1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
