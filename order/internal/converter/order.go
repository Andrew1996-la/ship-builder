package converter

import (
	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
)

func ToDTO(order model.Order) *orderv1.OrderDto {
	var shieldUUID orderv1.OptNilUUID
	if item, ok := order.ItemByType(model.PartTypeShield); ok {
		shieldUUID = orderv1.NewOptNilUUID(item.PartUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if item, ok := order.ItemByType(model.PartTypeWeapon); ok {
		weaponUUID = orderv1.NewOptNilUUID(item.PartUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		pm := orderv1.PaymentMethod(*order.PaymentMethod)
		paymentMethod = orderv1.NewOptNilPaymentMethod(pm)
	}

	dto := &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}

	if item, ok := order.ItemByType(model.PartTypeHull); ok {
		dto.HullUUID = item.PartUUID
	}
	if item, ok := order.ItemByType(model.PartTypeEngine); ok {
		dto.EngineUUID = item.PartUUID
	}

	return dto
}
