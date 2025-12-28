package schemas

import "time"

type Clients struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	Name      string
	Status    string    `gorm:"type:enum('active','inactive');default:'active'"`
	Telefone  string    `gorm:"type:varchar(20)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTIme"`
	Address   Address   `gorm:"type:json;serializer:json"`
}

type Address struct {
	Street       string // Rua
	Number       string // Número da casa
	Complement   string // Complemento
	Neighborhood string // Bairro
	City         string // Cidade
	State        string // Estado
	ZipCode      string //
}
