package model

import "github.com/google/uuid"

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

func (m PaymentMethod) IsValid() bool {
	switch m {
	case PaymentMethodCard, PaymentMethodSBP,
		PaymentMethodCreditCard, PaymentMethodInvestorMoney:
		return true
	default:
		return false
	}
}

type PayRequest struct {
	OrderUUID     uuid.UUID
	PaymentMethod PaymentMethod
}

type Payment struct {
	OrderUUID       uuid.UUID
	TransactionUUID uuid.UUID
	PaymentMethod   PaymentMethod
}
