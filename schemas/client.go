package schemas

import "time"

type Clients struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	StoreID   uint64    `gorm:"not null"`
	UserID    uint64    `gorm:"not null"`
	Name      string    `gorm:"type:varchar(255)"`
	Status    string    `gorm:"type:enum('active','inactive');default:'active'"`
	Role      UserRole  `json:"role" gorm:"type:enum('admin','user','client');default:'client'"`
	Phone     string    `gorm:"type:varchar(20)"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTIme"`
	Address   Address   `gorm:"type:json;serializer:json"`
}

type Address struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	ClientID     uint64 `gorm:"not null"`
	Street       string `gorm:"type:varchar(255)"` // Rua
	Number       string `gorm:"type:varchar(10)"`  // Número da casa
	Complement   string `gorm:"type:varchar(255)"` // Complemento
	Neighborhood string `gorm:"type:varchar(255)"` // Bairro
	City         string `gorm:"type:varchar(255)"` // Cidade
	State        string `gorm:"type:varchar(2)"`   // Estado
	ZipCode      string `gorm:"type:varchar(10)"`  // CEP
}
