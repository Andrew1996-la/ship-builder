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

	items := make([]record.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, record.OrderItem{
			UUID:      item.UUID,
			OrderUUID: item.OrderUUID,
			PartUUID:  item.PartUUID,
			PartType:  item.PartType.String(),
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
		})
	}

	return record.Order{
		OrderUUID:       order.OrderUUID,
		Items:           items,
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

	items := make([]model.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, model.OrderItem{
			UUID:      item.UUID,
			OrderUUID: item.OrderUUID,
			PartUUID:  item.PartUUID,
			PartType:  model.PartType(item.PartType),
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
		})
	}

	return model.Order{
		OrderUUID:       order.OrderUUID,
		Items:           items,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}
