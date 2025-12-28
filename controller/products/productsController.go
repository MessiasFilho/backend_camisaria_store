package controller

import (
	"backend_camisaria_store/config"
	"backend_camisaria_store/schemas"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	req := CreateProductRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Erro ao processar dados da requisição",
			"details": err.Error(),
		})
	}

	// Validar dados
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
	}

	// Mapear request para schema
	product := schemas.Products{
		SKU:              req.SKU,
		Name:             req.Name,
		Description:      req.Description,
		Categorys:        req.Categorys,
		Size:             req.Size,
		Color:            req.Color,
		Material:         req.Material,
		Gender:           req.Gender,
		Price:            req.Price,
		PromotionalPrice: req.PromotionalPrice,
		StockQuantity:    req.StockQuantity,
		MinStock:         req.MinStock,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		Tags:             req.Tags,
		SEODescription:   req.SEODescription,
		SEOKeywords:      req.SEOKeywords,
		IsActive:         true, // Produtos novos são ativos por padrão
		IsPromotional:    false,
	}

	// Salvar no banco
	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao salvar produto",
			"details": err.Error(),
		})
	}

	// Retornar resposta de sucesso
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Produto criado com sucesso",
	})
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	// Converter string para uint
	productID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	var product schemas.Products
	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Produto não encontrado",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product": ProductResponse{
			ID:               product.ID,
			SKU:              product.SKU,
			Name:             product.Name,
			Description:      product.Description,
			Categorys:        product.Categorys,
			Size:             product.Size,
			Color:            product.Color,
			Material:         product.Material,
			Gender:           product.Gender,
			Price:            product.Price,
			PromotionalPrice: product.PromotionalPrice,
			StockQuantity:    product.StockQuantity,
			MinStock:         product.MinStock,
			Weight:           product.Weight,
			Dimensions:       product.Dimensions,
			IsActive:         product.IsActive,
			IsPromotional:    product.IsPromotional,
			Tags:             product.Tags,
			SEODescription:   product.SEODescription,
			SEOKeywords:      product.SEOKeywords,
		},
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	// Converter string para uint
	productID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	req := UpdateProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Erro ao processar dados da requisição",
			"details": err.Error(),
		})
	}

	// Validar dados
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
	}

	// Verificar se produto existe
	var product schemas.Products
	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Produto não encontrado",
		})
	}

	// Mapear campos atualizáveis
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Categorys != nil {
		updates["categorys"] = *req.Categorys
	}
	if req.Subcategory != nil {
		updates["subcategory"] = *req.Subcategory
	}
	if req.Brand != nil {
		updates["brand"] = *req.Brand
	}
	if req.Size != nil {
		updates["size"] = *req.Size
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.Material != nil {
		updates["material"] = *req.Material
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.PromotionalPrice != nil {
		updates["promotional_price"] = req.PromotionalPrice
	}
	if req.StockQuantity != nil {
		updates["stock_quantity"] = *req.StockQuantity
	}
	if req.MinStock != nil {
		updates["min_stock"] = *req.MinStock
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.Dimensions != nil {
		updates["dimensions"] = *req.Dimensions
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsPromotional != nil {
		updates["is_promotional"] = *req.IsPromotional
	}
	if req.Tags != nil {
		updates["tags"] = *req.Tags
	}
	if req.SEODescription != nil {
		updates["seo_description"] = *req.SEODescription
	}
	if req.SEOKeywords != nil {
		updates["seo_keywords"] = *req.SEOKeywords
	}

	// Atualizar no banco
	if err := config.DB.Model(&product).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao atualizar produto",
			"details": err.Error(),
		})
	}

	// Buscar produto atualizado
	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar produto atualizado",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Produto atualizado com sucesso",
		"product": ProductResponse{
			ID:               product.ID,
			SKU:              product.SKU,
			Name:             product.Name,
			Description:      product.Description,
			Categorys:        product.Categorys,
			Size:             product.Size,
			Color:            product.Color,
			Material:         product.Material,
			Gender:           product.Gender,
			Price:            product.Price,
			PromotionalPrice: product.PromotionalPrice,
			StockQuantity:    product.StockQuantity,
			MinStock:         product.MinStock,
			Weight:           product.Weight,
			Dimensions:       product.Dimensions,
			IsActive:         product.IsActive,
			IsPromotional:    product.IsPromotional,
			Tags:             product.Tags,
			SEODescription:   product.SEODescription,
			SEOKeywords:      product.SEOKeywords,
		},
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	// Converter string para uint
	productID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	// Verificar se produto existe
	var product schemas.Products
	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Produto não encontrado",
		})
	}

	// Soft delete (desativar ao invés de deletar)
	if err := config.DB.Model(&product).Update("is_active", false).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao deletar produto",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Produto deletado com sucesso",
	})
}

func ListProducts(c *fiber.Ctx) error {
	// Parâmetros de paginação
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Filtros
	filters := ProductFilter{}
	if category := c.Query("category"); category != "" {
		cat := schemas.Category(category)
		filters.Category = &cat
	}
	if subcategory := c.Query("subcategory"); subcategory != "" {
		filters.Subcategory = &subcategory
	}
	if brand := c.Query("brand"); brand != "" {
		filters.Brand = &brand
	}
	if size := c.Query("size"); size != "" {
		filters.Size = &size
	}
	if color := c.Query("color"); color != "" {
		filters.Color = &color
	}
	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	// Query base
	query := config.DB.Model(&schemas.Products{})

	// Aplicar filtros
	if filters.Category != nil {
		query = query.Where("categorys = ?", *filters.Category)
	}
	if filters.Subcategory != nil {
		query = query.Where("subcategory LIKE ?", "%"+*filters.Subcategory+"%")
	}
	if filters.Brand != nil {
		query = query.Where("brand LIKE ?", "%"+*filters.Brand+"%")
	}
	if filters.Size != nil {
		query = query.Where("size = ?", *filters.Size)
	}
	if filters.Color != nil {
		query = query.Where("color LIKE ?", "%"+*filters.Color+"%")
	}
	if filters.Search != nil {
		query = query.Where("name LIKE ? OR description LIKE ? OR sku LIKE ?",
			"%"+*filters.Search+"%", "%"+*filters.Search+"%", "%"+*filters.Search+"%")
	}

	// Contar total de registros
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao contar produtos",
		})
	}

	// Buscar produtos
	var products []schemas.Products
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar produtos",
		})
	}

	// Converter para response
	var productResponses []ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, ProductResponse{
			ID:               product.ID,
			SKU:              product.SKU,
			Name:             product.Name,
			Description:      product.Description,
			Categorys:        product.Categorys,
			Size:             product.Size,
			Color:            product.Color,
			Material:         product.Material,
			Gender:           product.Gender,
			Price:            product.Price,
			PromotionalPrice: product.PromotionalPrice,
			StockQuantity:    product.StockQuantity,
			MinStock:         product.MinStock,
			Weight:           product.Weight,
			Dimensions:       product.Dimensions,
			IsActive:         product.IsActive,
			IsPromotional:    product.IsPromotional,
			Tags:             product.Tags,
			SEODescription:   product.SEODescription,
			SEOKeywords:      product.SEOKeywords,
		})
	}

	// Calcular total de páginas
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.Status(fiber.StatusOK).JSON(ProductListResponse{
		Products: productResponses,
		Total:    total,
		Page:     page,
		Limit:    limit,
		Pages:    totalPages,
	})
}
