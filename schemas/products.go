package schemas

import (
	"time"
)

type Category string

const (
	Masculino   Category = "masculino"
	Feminino    Category = "Femenino"
	Fardamentos Category = "Fardamentos"
)

type Products struct {
	ID               uint64 `gorm:"primaryKey;autoIncrement"`
	SKU              string `gorm:"unique;not null"`
	Name             string `gorm:"not null"`
	Description      string
	Categorys        Category       `gorm:"not null;default:'masculino'"` // feminino, masculino, fardamentos
	Size             string         `gorm:"not null"`                     // PP, P, M, G, GG, XG, etc.
	Color            string         `gorm:"not null"`
	Material         string         // algodão, poliéster, etc.
	Gender           string         // M, F, U (unisex)
	Price            float64        `gorm:"type:decimal(10,2);not null"`
	PromotionalPrice *float64       `gorm:"type:decimal(10,2)"`
	StockQuantity    int            `gorm:"default:0"`
	MinStock         int            `gorm:"default:0"`
	Weight           float64        `gorm:"type:decimal(5,2)"` // em kg
	Dimensions       string         // comprimento x largura x altura
	Images           []string `gorm:"type:json"`
	IsActive         bool           `gorm:"default:true"`
	IsPromotional    bool           `gorm:"default:false"`
	Tags             string         // tags separadas por vírgula
	SEODescription   string
	SEOKeywords      string
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

type ProductImage struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	ProductID uint64
	URL       string `gorm:"not null"`
	AltText   string
	IsMain    bool      `gorm:"default:false"`
	Order     int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
