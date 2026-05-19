package schemas

import "time"

type Instance struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Integration string    `gorm:"not null"`
	QrCode      bool      `gorm:"type:boolean;not null"`
	Status      string    `gorm:"type:enum('active','inactive');default:'active'"`
	Hash        string    `gorm:"type:varchar(50)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
