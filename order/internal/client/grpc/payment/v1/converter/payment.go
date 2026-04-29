package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

func ToProtoPaymentMethod(paymentMethod model.PaymentMethod) paymentv1.PaymentMethod {
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

func ToModelTransactionUUID(resp *paymentv1.PayOrderResponse) (uuid.UUID, error) {
	transactionUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, fmt.Errorf("разобрать UUID транзакции: %w", err)
	}

	return transactionUUID, nil
}
