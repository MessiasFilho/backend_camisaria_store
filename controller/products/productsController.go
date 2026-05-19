package controller

import (
	"backend_camisaria_store/config"
	"backend_camisaria_store/schemas"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func parsePagination(c *fiber.Ctx) (page, limit, offset int) {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err = strconv.Atoi(c.Query("limit", "12"))
	if err != nil || limit < 1 {
		limit = 12
	}
	if limit > 100 {
		limit = 100
	}
	offset = (page - 1) * limit
	return page, limit, offset
}

func applyProductFilters(query *gorm.DB, filters ProductFilter) *gorm.DB {
	if filters.Category != nil {
		query = query.Where("LOWER(categorys) = ?", strings.ToLower(string(*filters.Category)))
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Active != nil {
		query = query.Where("is_active = ?", *filters.Active)
	}
	if filters.Search != nil && strings.TrimSpace(*filters.Search) != "" {
		term := "%" + strings.TrimSpace(*filters.Search) + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR sku LIKE ?", term, term, term)
	}
	return query
}

func listProductsWithFilters(c *fiber.Ctx, filters ProductFilter) error {
	page, limit, offset := parsePagination(c)

	query := config.DB.Model(&schemas.Products{})
	query = applyProductFilters(query, filters)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao contar produtos",
		})
	}

	var products []schemas.Products
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar produtos",
		})
	}

	responses := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		responses = append(responses, toProductResponse(p))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return c.Status(fiber.StatusOK).JSON(ProductListResponse{
		Products: responses,
		Total:    total,
		Page:     page,
		Limit:    limit,
		Pages:    totalPages,
	})
}

// ListCategoriesSummary — totais por categoria para o aside do admin.
func ListCategoriesSummary(c *fiber.Ctx) error {
	type row struct {
		Category string `gorm:"column:category"`
		Count    int64  `gorm:"column:count"`
	}
	var rows []row

	err := config.DB.Raw(`
		SELECT categorys AS category, COUNT(*) AS count
		FROM products
		WHERE is_active = ?
		GROUP BY categorys
	`, true).Scan(&rows).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao buscar resumo de categorias",
			"details": err.Error(),
		})
	}

	counts := map[schemas.Category]int64{
		schemas.Masculino:   0,
		schemas.Feminino:    0,
		schemas.Fardamentos: 0,
	}
	var total int64
	for _, r := range rows {
		cat := schemas.Category(strings.ToLower(strings.TrimSpace(r.Category)))
		if isValidCategory(cat) {
			counts[cat] += r.Count
			total += r.Count
		}
	}

	categories := []CategoryCountItem{
		{Category: schemas.Masculino, Count: counts[schemas.Masculino]},
		{Category: schemas.Feminino, Count: counts[schemas.Feminino]},
		{Category: schemas.Fardamentos, Count: counts[schemas.Fardamentos]},
	}

	return c.Status(fiber.StatusOK).JSON(CategoriesSummaryResponse{
		Categories: categories,
		Total:      total,
	})
}

// ListProductsByCategory — lista paginada filtrada pela categoria (path).
func ListProductsByCategory(c *fiber.Ctx) error {
	categoryParam := c.Params("category")
	cat := schemas.Category(categoryParam)
	if !isValidCategory(cat) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Categoria inválida. Use: masculino, feminino ou fardamentos",
		})
	}

	filters := ProductFilter{Category: &cat}

	if status := c.Query("status"); status != "" {
		st := schemas.ProductStatus(status)
		filters.Status = &st
	}
	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}
	active := true
	if activeStr := c.Query("active"); activeStr != "" {
		active = activeStr == "true" || activeStr == "1"
	}
	filters.Active = &active

	return listProductsWithFilters(c, filters)
}

// ListProducts — admin: lista paginada com filtros (category, status, search, active).
func ListProducts(c *fiber.Ctx) error {
	filters := ProductFilter{}

	if category := c.Query("category"); category != "" {
		cat := schemas.Category(category)
		filters.Category = &cat
	}
	if status := c.Query("status"); status != "" {
		st := schemas.ProductStatus(status)
		filters.Status = &st
	}
	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}
	if activeStr := c.Query("active"); activeStr != "" {
		active := activeStr == "true" || activeStr == "1"
		filters.Active = &active
	} else {
		// Admin vê apenas registros ativos por padrão (soft delete mantém is_active=false)
		active := true
		filters.Active = &active
	}

	return listProductsWithFilters(c, filters)
}

// ListPublishedProducts — loja pública: apenas publicados e ativos.
func ListPublishedProducts(c *fiber.Ctx) error {
	published := schemas.ProductStatusPublished
	active := true
	filters := ProductFilter{
		Status: &published,
		Active: &active,
	}

	if category := c.Query("category"); category != "" {
		cat := schemas.Category(category)
		if isValidCategory(cat) {
			filters.Category = &cat
		}
	}
	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	return listProductsWithFilters(c, filters)
}

func CreateProduct(c *fiber.Ctx) error {
	req := CreateProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Erro ao processar dados da requisição",
			"details": err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
	}

	status := req.Status
	if status == "" {
		status = schemas.ProductStatusDraft
	}

	isPromotional := req.PromotionalPrice != nil && *req.PromotionalPrice > 0

	product := schemas.Products{
		SKU:              strings.TrimSpace(req.SKU),
		Name:             strings.TrimSpace(req.Name),
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
		Status:           status,
		IsActive:         true,
		IsPromotional:    isPromotional,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao salvar produto",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Produto criado com sucesso",
		"product": toProductResponse(product),
	})
}

func GetProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var product schemas.Products
	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Produto não encontrado"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product": toProductResponse(product),
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	req := UpdateProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Erro ao processar dados da requisição",
			"details": err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
	}

	var product schemas.Products
	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Produto não encontrado"})
	}

	updates := make(map[string]interface{})

	if req.SKU != nil {
		updates["sku"] = strings.TrimSpace(*req.SKU)
	}
	if req.Name != nil {
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Categorys != nil {
		updates["categorys"] = *req.Categorys
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
		if *req.PromotionalPrice > 0 {
			updates["is_promotional"] = true
		} else {
			updates["is_promotional"] = false
		}
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
	if req.Status != nil {
		updates["status"] = *req.Status
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

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nenhum campo para atualizar",
		})
	}

	if err := config.DB.Model(&product).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao atualizar produto",
			"details": err.Error(),
		})
	}

	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar produto atualizado",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Produto atualizado com sucesso",
		"product": toProductResponse(product),
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var product schemas.Products
	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Produto não encontrado"})
	}

	if err := config.DB.Model(&product).Updates(map[string]interface{}{
		"is_active": false,
		"status":    schemas.ProductStatusDraft,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao excluir produto",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Produto removido com sucesso",
	})
}
