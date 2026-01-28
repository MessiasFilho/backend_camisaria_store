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
	ID               uint64    `gorm:"primaryKey;autoIncrement"`
	SKU              string    `gorm:"type:varchar(100);uniqueIndex:uni_products_sku,size:100;not null"`
	Name             string    `gorm:"type:varchar(255);not null"`
	Description      string    `gorm:"type:text"`
	Categorys        Category  `gorm:"not null;default:'masculino'"` // feminino, masculino, fardamentos
	Size             string    `gorm:"type:varchar(10);not null"`    // PP, P, M, G, GG, XG, etc.
	Color            string    `gorm:"type:varchar(50);not null"`
	Material         string    `gorm:"type:varchar(100)"` // algodão, poliéster, etc.
	Gender           string    `gorm:"type:varchar(1)"`   // M, F, U (unisex)
	Price            float64   `gorm:"type:decimal(10,2);not null"`
	PromotionalPrice *float64  `gorm:"type:decimal(10,2)"`
	StockQuantity    int       `gorm:"default:0"`
	MinStock         int       `gorm:"default:0"`
	Weight           float64   `gorm:"type:decimal(5,2)"` // em kg
	Dimensions       string    `gorm:"type:varchar(100)"` // comprimento x largura x altura
	Images           []byte    `gorm:"type:json"`
	IsActive         bool      `gorm:"default:true"`
	IsPromotional    bool      `gorm:"default:false"`
	Tags             string    `gorm:"type:text"` // tags separadas por vírgula
	SEODescription   string    `gorm:"type:text"`
	SEOKeywords      string    `gorm:"type:text"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}
