package schemas

import "time"

type StatusPayment string

const (
	PendingPayment StatusPayment = "pending"
	PaidPayment    StatusPayment = "paid"
	FailedPayment  StatusPayment = "failed"
)

type DeliveryType string

const (
	PickupDelivery   DeliveryType = "pickup"
	DeliveryDelivery DeliveryType = "delivery"
)

type Orders struct {
	ID            uint64        `gorm:"primaryKey;autoIncrement"`
	ClientID      uint64        `gorm:"not null"`
	OrderNumber   string        `gorm:"unique;not null"`
	Value         float64       `gorm:"type:decimal(10,2);not null"`
	Originalvalue float64       `gorm:"type:decimal(10,2);not null"`
	StatusPayment StatusPayment `gorm:"type:enum('pending','paid','failed');default:'pending'"`
	DeliveryType  DeliveryType  `gorm:"type:enum('pickup','delivery');default:'pickup'"`
	CreatedAt     time.Time     `gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime"`
}
