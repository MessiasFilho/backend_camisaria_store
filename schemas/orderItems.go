package schemas

import "time"

type OrderItems struct {
	ID            uint64  `gorm:"primaryKey;autoIncrement"`
	OrderID       uint64  `gorm:"not null"`
	ProductID     uint64  `gorm:"not null"`
	Quantity      int     `gorm:"not null"`
	Price         float64 `gorm:"type:decimal(10,2);not null"`
	OriginalPrice float64 `gorm:"type:decimal(10,2);not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
