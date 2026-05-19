package controller

import (
	"errors"
	"strings"
	"time"

	"backend_camisaria_store/schemas"
	"backend_camisaria_store/service/minio"
)

type CreateProductRequest struct {
	SKU              string                `json:"sku"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	Categorys        schemas.Category      `json:"categorys"`
	Size             string                `json:"size"`
	Color            string                `json:"color"`
	Material         string                `json:"material"`
	Gender           string                `json:"gender"`
	Price            float64               `json:"price"`
	PromotionalPrice *float64              `json:"promotional_price,omitempty"`
	StockQuantity    int                   `json:"stock_quantity"`
	MinStock         int                   `json:"min_stock"`
	Weight           float64               `json:"weight"`
	Dimensions       string                `json:"dimensions"`
	Tags             string                `json:"tags"`
	SEODescription   string                `json:"seo_description"`
	SEOKeywords      string                `json:"seo_keywords"`
	Status           schemas.ProductStatus `json:"status"`
}

type UpdateProductRequest struct {
	SKU              *string                `json:"sku,omitempty"`
	Name             *string                `json:"name,omitempty"`
	Description      *string                `json:"description,omitempty"`
	Categorys        *schemas.Category      `json:"categorys,omitempty"`
	Size             *string                `json:"size,omitempty"`
	Color            *string                `json:"color,omitempty"`
	Material         *string                `json:"material,omitempty"`
	Gender           *string                `json:"gender,omitempty"`
	Price            *float64               `json:"price,omitempty"`
	PromotionalPrice *float64               `json:"promotional_price,omitempty"`
	StockQuantity    *int                   `json:"stock_quantity,omitempty"`
	MinStock         *int                   `json:"min_stock,omitempty"`
	Weight           *float64               `json:"weight,omitempty"`
	Dimensions       *string                `json:"dimensions,omitempty"`
	IsActive         *bool                  `json:"is_active,omitempty"`
	IsPromotional    *bool                  `json:"is_promotional,omitempty"`
	Status           *schemas.ProductStatus `json:"status,omitempty"`
	Tags             *string                `json:"tags,omitempty"`
	SEODescription   *string                `json:"seo_description,omitempty"`
	SEOKeywords      *string                `json:"seo_keywords,omitempty"`
}

type ProductResponse struct {
	ID               uint64                `json:"id"`
	SKU              string                `json:"sku"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	Categorys        schemas.Category      `json:"categorys"`
	Size             string                `json:"size"`
	Color            string                `json:"color"`
	Material         string                `json:"material"`
	Gender           string                `json:"gender"`
	Price            float64               `json:"price"`
	PromotionalPrice *float64              `json:"promotional_price,omitempty"`
	StockQuantity    int                   `json:"stock_quantity"`
	MinStock         int                   `json:"min_stock"`
	Weight           float64               `json:"weight"`
	Dimensions       string                `json:"dimensions"`
	Images           []string              `json:"images"`
	Status           schemas.ProductStatus `json:"status"`
	IsActive         bool                  `json:"is_active"`
	IsPromotional    bool                  `json:"is_promotional"`
	Tags             string                `json:"tags"`
	SEODescription   string                `json:"seo_description"`
	SEOKeywords      string                `json:"seo_keywords"`
	CreatedAt        string                `json:"created_at"`
	UpdatedAt        string                `json:"updated_at"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Pages    int               `json:"pages"`
}

type CategoryCountItem struct {
	Category schemas.Category `json:"category"`
	Count    int64            `json:"count"`
}

type CategoriesSummaryResponse struct {
	Categories []CategoryCountItem `json:"categories"`
	Total      int64               `json:"total"`
}

type ProductFilter struct {
	Category *schemas.Category
	Status   *schemas.ProductStatus
	Search   *string
	Active   *bool
}

func toProductResponse(p schemas.Products) ProductResponse {
	return ProductResponse{
		ID:               p.ID,
		SKU:              p.SKU,
		Name:             p.Name,
		Description:      p.Description,
		Categorys:        p.Categorys,
		Size:             p.Size,
		Color:            p.Color,
		Material:         p.Material,
		Gender:           p.Gender,
		Price:            p.Price,
		PromotionalPrice: p.PromotionalPrice,
		StockQuantity:    p.StockQuantity,
		MinStock:         p.MinStock,
		Weight:           p.Weight,
		Dimensions:       p.Dimensions,
		Images:           minio.JsonToStringSlice(p.Images),
		Status:           p.Status,
		IsActive:         p.IsActive,
		IsPromotional:    p.IsPromotional,
		Tags:             p.Tags,
		SEODescription:   p.SEODescription,
		SEOKeywords:      p.SEOKeywords,
		CreatedAt:        p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        p.UpdatedAt.Format(time.RFC3339),
	}
}

func isValidCategory(c schemas.Category) bool {
	return c == schemas.Masculino || c == schemas.Feminino || c == schemas.Fardamentos
}

func isValidStatus(s schemas.ProductStatus) bool {
	return s == schemas.ProductStatusDraft || s == schemas.ProductStatusPublished
}

// genderForCategory mantém gênero alinhado à categoria do catálogo.
func genderForCategory(cat schemas.Category) string {
	switch cat {
	case schemas.Masculino:
		return "M"
	case schemas.Feminino:
		return "F"
	case schemas.Fardamentos:
		return "U"
	default:
		return "U"
	}
}

func (req *CreateProductRequest) Validate() error {
	var errs []string

	if strings.TrimSpace(req.SKU) == "" {
		errs = append(errs, "SKU é obrigatório")
	} else if len(req.SKU) < 3 || len(req.SKU) > 50 {
		errs = append(errs, "SKU deve ter entre 3 e 50 caracteres")
	}

	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, "nome é obrigatório")
	} else if len(req.Name) < 3 || len(req.Name) > 255 {
		errs = append(errs, "nome deve ter entre 3 e 255 caracteres")
	}

	if !isValidCategory(req.Categorys) {
		errs = append(errs, "categoria deve ser masculino, feminino ou fardamentos")
	}

	if strings.TrimSpace(req.Size) == "" {
		errs = append(errs, "tamanho é obrigatório")
	}

	if strings.TrimSpace(req.Color) == "" {
		errs = append(errs, "cor é obrigatória")
	}

	req.Gender = genderForCategory(req.Categorys)

	if req.Price <= 0 {
		errs = append(errs, "preço deve ser maior que zero")
	}

	if req.PromotionalPrice != nil && *req.PromotionalPrice < 0 {
		errs = append(errs, "preço promocional não pode ser negativo")
	}

	if req.StockQuantity < 0 {
		errs = append(errs, "quantidade em estoque não pode ser negativa")
	}

	if req.MinStock < 0 {
		errs = append(errs, "estoque mínimo não pode ser negativo")
	}

	if req.Weight < 0 {
		errs = append(errs, "peso não pode ser negativo")
	}

	if req.Status != "" && !isValidStatus(req.Status) {
		errs = append(errs, "status deve ser draft ou published")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func (req *UpdateProductRequest) Validate() error {
	var errs []string

	if req.SKU != nil {
		if strings.TrimSpace(*req.SKU) == "" {
			errs = append(errs, "SKU não pode ser vazio")
		} else if len(*req.SKU) < 3 || len(*req.SKU) > 50 {
			errs = append(errs, "SKU deve ter entre 3 e 50 caracteres")
		}
	}

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			errs = append(errs, "nome não pode ser vazio")
		} else if len(*req.Name) < 3 || len(*req.Name) > 255 {
			errs = append(errs, "nome deve ter entre 3 e 255 caracteres")
		}
	}

	if req.Categorys != nil && !isValidCategory(*req.Categorys) {
		errs = append(errs, "categoria deve ser masculino, feminino ou fardamentos")
	}

	if req.Categorys != nil {
		g := genderForCategory(*req.Categorys)
		req.Gender = &g
	}

	if req.Price != nil && *req.Price <= 0 {
		errs = append(errs, "preço deve ser maior que zero")
	}

	if req.PromotionalPrice != nil && *req.PromotionalPrice < 0 {
		errs = append(errs, "preço promocional não pode ser negativo")
	}

	if req.StockQuantity != nil && *req.StockQuantity < 0 {
		errs = append(errs, "quantidade em estoque não pode ser negativa")
	}

	if req.MinStock != nil && *req.MinStock < 0 {
		errs = append(errs, "estoque mínimo não pode ser negativo")
	}

	if req.Weight != nil && *req.Weight < 0 {
		errs = append(errs, "peso não pode ser negativo")
	}

	if req.Status != nil && !isValidStatus(*req.Status) {
		errs = append(errs, "status deve ser draft ou published")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}
