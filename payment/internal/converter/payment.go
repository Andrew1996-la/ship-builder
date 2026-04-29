package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/payment/internal/model"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

func ToModelPaymentMethod(paymentMethod paymentv1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

func ToModelPayRequest(req *paymentv1.PayOrderRequest) (model.PayRequest, error) {
	orderUuid, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return model.PayRequest{}, fmt.Errorf("разобрать UUID заказа: %w", err)
	}

	return model.PayRequest{
		OrderUUID:     orderUuid,
		PaymentMethod: ToModelPaymentMethod(req.GetPaymentMethod()),
	}, nil
}

func ToProtoPayOrderResponse(payment model.Payment) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: payment.TransactionUUID.String(),
	}
}
