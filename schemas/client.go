package schemas

import "time"

type Clients struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	StoreID   uint64 `gorm:"not null"`
	UserID    uint64 `gorm:"not null"`
	Name      string
	Status    string    `gorm:"type:enum('active','inactive');default:'active'"`
	Role      UserRole  `json:"role" gorm:"type:enum('admin','user','client');default:'client'"`
	Phone     string    `gorm:"type:varchar(20)"`
	Email     string    `gorm:"type:varchar(255);unique"`
	Password  string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTIme"`
	Address   Address   `gorm:"type:json;serializer:json"`
}

type Address struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	ClientID     uint64 `gorm:"not null"`
	Street       string // Rua
	Number       string // Número da casa
	Complement   string // Complemento
	Neighborhood string // Bairro
	City         string // Cidade
	State        string // Estado
	ZipCode      string //
}
