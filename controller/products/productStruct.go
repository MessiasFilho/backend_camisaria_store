package controller

import (
	"errors"
	"strings"

	"backend_camisaria_store/schemas"
)

// CreateProductRequest representa a requisição para criar um produto
type CreateProductRequest struct {
	SKU              string           `json:"sku" validate:"required,min=3,max=50"`
	Name             string           `json:"name" validate:"required,min=3,max=255"`
	Description      string           `json:"description"`
	Categorys        schemas.Category `json:"categorys" validate:"required,oneof=masculino feminino fardamentos"`
	Subcategory      string           `json:"subcategory"`
	Brand            string           `json:"brand"`
	Size             string           `json:"size" validate:"required"`
	Color            string           `json:"color" validate:"required"`
	Material         string           `json:"material"`
	Gender           string           `json:"gender" validate:"omitempty,oneof=M F U"`
	Price            float64          `json:"price" validate:"required,min=0.01"`
	PromotionalPrice *float64         `json:"promotional_price,omitempty" validate:"omitempty,min=0"`
	StockQuantity    int              `json:"stock_quantity" validate:"min=0"`
	MinStock         int              `json:"min_stock" validate:"min=0"`
	Weight           float64          `json:"weight" validate:"min=0"`
	Dimensions       string           `json:"dimensions"`
	Tags             string           `json:"tags"`
	SEODescription   string           `json:"seo_description"`
	SEOKeywords      string           `json:"seo_keywords"`
}

// UpdateProductRequest representa a requisição para atualizar um produto
type UpdateProductRequest struct {
	Name             *string           `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description      *string           `json:"description,omitempty"`
	Categorys        *schemas.Category `json:"categorys,omitempty" validate:"omitempty,oneof=masculino feminino fardamentos"`
	Subcategory      *string           `json:"subcategory,omitempty"`
	Brand            *string           `json:"brand,omitempty"`
	Size             *string           `json:"size,omitempty"`
	Color            *string           `json:"color,omitempty"`
	Material         *string           `json:"material,omitempty"`
	Gender           *string           `json:"gender,omitempty" validate:"omitempty,oneof=M F U"`
	Price            *float64          `json:"price,omitempty" validate:"omitempty,min=0.01"`
	PromotionalPrice *float64          `json:"promotional_price,omitempty" validate:"omitempty,min=0"`
	StockQuantity    *int              `json:"stock_quantity,omitempty" validate:"omitempty,min=0"`
	MinStock         *int              `json:"min_stock,omitempty" validate:"omitempty,min=0"`
	Weight           *float64          `json:"weight,omitempty" validate:"min=0"`
	Dimensions       *string           `json:"dimensions,omitempty"`
	IsActive         *bool             `json:"is_active,omitempty"`
	IsPromotional    *bool             `json:"is_promotional,omitempty"`
	Tags             *string           `json:"tags,omitempty"`
	SEODescription   *string           `json:"seo_description,omitempty"`
	SEOKeywords      *string           `json:"seo_keywords,omitempty"`
}

// ProductResponse representa a resposta da API para produtos
type ProductResponse struct {
	ID               uint64           `json:"id"`
	SKU              string           `json:"sku"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Categorys        schemas.Category `json:"categorys"`
	Subcategory      string           `json:"subcategory"`
	Brand            string           `json:"brand"`
	Size             string           `json:"size"`
	Color            string           `json:"color"`
	Material         string           `json:"material"`
	Gender           string           `json:"gender"`
	Price            float64          `json:"price"`
	PromotionalPrice *float64         `json:"promotional_price,omitempty"`
	StockQuantity    int              `json:"stock_quantity"`
	MinStock         int              `json:"min_stock"`
	Weight           float64          `json:"weight"`
	Dimensions       string           `json:"dimensions"`
	IsActive         bool             `json:"is_active"`
	IsPromotional    bool             `json:"is_promotional"`
	Tags             string           `json:"tags"`
	SEODescription   string           `json:"seo_description"`
	SEOKeywords      string           `json:"seo_keywords"`
	CreatedAt        string           `json:"created_at"`
	UpdatedAt        string           `json:"updated_at"`
}

// ProductListResponse representa a resposta paginada para listagem de produtos
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Pages    int               `json:"pages"`
}

// ProductFilter representa os filtros para busca de produtos
type ProductFilter struct {
	Category    *schemas.Category `json:"category,omitempty"`
	Subcategory *string           `json:"subcategory,omitempty"`
	Brand       *string           `json:"brand,omitempty"`
	Size        *string           `json:"size,omitempty"`
	Color       *string           `json:"color,omitempty"`
	MinPrice    *float64          `json:"min_price,omitempty"`
	MaxPrice    *float64          `json:"max_price,omitempty"`
	IsActive    *bool             `json:"is_active,omitempty"`
	Search      *string           `json:"search,omitempty"`
}

// PaginationParams representa os parâmetros de paginação
type PaginationParams struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

// Validate realiza validação personalizada dos campos do CreateProductRequest
func (req *CreateProductRequest) Validate() error {
	var errs []string

	// Validação do SKU
	if strings.TrimSpace(req.SKU) == "" {
		errs = append(errs, "SKU é obrigatório")
	} else if len(req.SKU) < 3 || len(req.SKU) > 50 {
		errs = append(errs, "SKU deve ter entre 3 e 50 caracteres")
	}

	// Validação do Nome
	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, "nome é obrigatório")
	} else if len(req.Name) < 3 || len(req.Name) > 255 {
		errs = append(errs, "nome deve ter entre 3 e 255 caracteres")
	}

	// Validação da Categoria
	if req.Categorys != schemas.Masculino && req.Categorys != schemas.Feminino && req.Categorys != schemas.Fardamentos {
		errs = append(errs, "categoria deve ser uma das seguintes: masculino, feminino, fardamentos")
	}

	// Validação do Tamanho
	if strings.TrimSpace(req.Size) == "" {
		errs = append(errs, "tamanho é obrigatório")
	}

	// Validação da Cor
	if strings.TrimSpace(req.Color) == "" {
		errs = append(errs, "cor é obrigatória")
	}

	// Validação do Gênero (opcional mas se informado deve ser válido)
	if req.Gender != "" && req.Gender != "M" && req.Gender != "F" && req.Gender != "U" {
		errs = append(errs, "gênero deve ser M, F ou U")
	}

	// Validação do Preço
	if req.Price <= 0 {
		errs = append(errs, "preço deve ser maior que zero")
	}

	// Validação do Preço Promocional (se informado)
	if req.PromotionalPrice != nil && *req.PromotionalPrice < 0 {
		errs = append(errs, "preço promocional não pode ser negativo")
	}

	// Validação do Estoque
	if req.StockQuantity < 0 {
		errs = append(errs, "quantidade em estoque não pode ser negativa")
	}

	if req.MinStock < 0 {
		errs = append(errs, "estoque mínimo não pode ser negativo")
	}

	// Validação do Peso
	if req.Weight < 0 {
		errs = append(errs, "peso não pode ser negativo")
	}

	// Se houver erros, retorna um erro composto
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// Validate realiza validação personalizada dos campos do UpdateProductRequest
func (req *UpdateProductRequest) Validate() error {
	var errs []string

	// Validação do Nome (se informado)
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			errs = append(errs, "nome não pode ser vazio")
		} else if len(*req.Name) < 3 || len(*req.Name) > 255 {
			errs = append(errs, "nome deve ter entre 3 e 255 caracteres")
		}
	}

	// Validação da Categoria (se informada)
	if req.Categorys != nil && *req.Categorys != schemas.Masculino && *req.Categorys != schemas.Feminino && *req.Categorys != schemas.Fardamentos {
		errs = append(errs, "categoria deve ser uma das seguintes: masculino, feminino, fardamentos")
	}

	// Validação do Gênero (se informado)
	if req.Gender != nil && *req.Gender != "M" && *req.Gender != "F" && *req.Gender != "U" {
		errs = append(errs, "gênero deve ser M, F ou U")
	}

	// Validação do Preço (se informado)
	if req.Price != nil && *req.Price <= 0 {
		errs = append(errs, "preço deve ser maior que zero")
	}

	// Validação do Preço Promocional (se informado)
	if req.PromotionalPrice != nil && *req.PromotionalPrice < 0 {
		errs = append(errs, "preço promocional não pode ser negativo")
	}

	// Validação do Estoque (se informado)
	if req.StockQuantity != nil && *req.StockQuantity < 0 {
		errs = append(errs, "quantidade em estoque não pode ser negativa")
	}

	if req.MinStock != nil && *req.MinStock < 0 {
		errs = append(errs, "estoque mínimo não pode ser negativo")
	}

	// Validação do Peso (se informado)
	if req.Weight != nil && *req.Weight < 0 {
		errs = append(errs, "peso não pode ser negativo")
	}

	// Se houver erros, retorna um erro composto
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
