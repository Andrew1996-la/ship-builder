package record

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderUUID       uuid.UUID
	Items           []OrderItem
	TotalPrice      int64 // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          string // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}

type OrderItem struct {
	UUID      uuid.UUID
	OrderUUID uuid.UUID
	PartUUID  uuid.UUID
	PartType  string
	Price     int64
	CreatedAt time.Time
}
