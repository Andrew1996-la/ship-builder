package converter

import (
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

func ToModelPayRequest(req *paymentv1.PayOrderRequest) model.PayRequest {
	return model.PayRequest{
		OrderUUID:     req.GetOrderUuid(),
		PaymentMethod: ToModelPaymentMethod(req.GetPaymentMethod()),
	}
}

func ToProtoPayOrderResponse(payment model.Payment) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: payment.TransactionUUID.String(),
	}
}
