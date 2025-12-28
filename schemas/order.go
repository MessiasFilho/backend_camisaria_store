package schemas

import "time"

type StatusOrder string

const (
	Pending    StatusOrder = "pending"
	Confirmed  StatusOrder = "confirmed"
	Processing StatusOrder = "processing"
	Shipped    StatusOrder = "shipped"
	Delivered  StatusOrder = "delivered"
	Cancelled  StatusOrder = "cancelled"
)

type Orders struct {
	ID          uint64      `gorm:"primaryKey;autoIncrement"`
	ClientID    uint64      `gorm:"not null"`
	OrderNumber string      `gorm:"unique;not null"`
	Status      StatusOrder `gorm:"type:enum('pending','confirmed','processing','shipped','delivered','cancelled');default:'pending'"`
	Total       float64     `gorm:"type:decimal(10,2);not null"`
	Notes       string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
