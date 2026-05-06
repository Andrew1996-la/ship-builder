package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type PartType string

const (
	PartTypeHull   PartType = "PART_TYPE_HULL"
	PartTypeEngine PartType = "PART_TYPE_ENGINE"
	PartTypeShield PartType = "PART_TYPE_SHIELD"
	PartTypeWeapon PartType = "PART_TYPE_WEAPON"
)

func (pt PartType) String() string {
	return string(pt)
}

type OrderItem struct {
	UUID      uuid.UUID
	OrderUUID uuid.UUID
	PartUUID  uuid.UUID
	PartType  PartType
	Price     int64
	CreatedAt time.Time
}

type Order struct {
	OrderUUID       uuid.UUID
	Items           []OrderItem
	TotalPrice      int64 // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

func (o Order) ItemByType(partType PartType) (OrderItem, bool) {
	for _, item := range o.Items {
		if item.PartType == partType {
			return item, true
		}
	}

	return OrderItem{}, false
}

type CreateOrderInfo struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

func (r *CreateOrderInfo) PartUUIDs() []uuid.UUID {
	uuids := []uuid.UUID{r.HullUUID, r.EngineUUID}
	if r.ShieldUUID != nil {
		uuids = append(uuids, *r.ShieldUUID)
	}
	if r.WeaponUUID != nil {
		uuids = append(uuids, *r.WeaponUUID)
	}
	return uuids
}

type PayOrderInfo struct {
	OrderUUID     uuid.UUID
	PaymentMethod PaymentMethod
}
