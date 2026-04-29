package converter

import (
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	"github.com/Andrew1996-la/ship-builder/order/internal/repository/record"
)

func ToRepoOrder(order model.Order) record.Order {
	var paymentMethod *string
	if order.PaymentMethod != nil {
		val := string(*order.PaymentMethod)
		paymentMethod = &val
	}

	return record.Order{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      order.ShieldUUID,
		WeaponUUID:      order.WeaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          string(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func ToModelOrder(order record.Order) model.Order {
	var paymentMethod *model.PaymentMethod
	if order.PaymentMethod != nil {
		val := model.PaymentMethod(*order.PaymentMethod)
		paymentMethod = &val
	}

	return model.Order{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      order.ShieldUUID,
		WeaponUUID:      order.WeaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}
