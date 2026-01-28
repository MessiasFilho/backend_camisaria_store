package controller

import (
	"backend_camisaria_store/config"
	"backend_camisaria_store/schemas"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateOrder(c *fiber.Ctx) error {

	req := CreateOrderRequest{}
	clientID := c.Locals("user_id").(uint64)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao processar dados da requisição",
		})
	}

	// Validar dados
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	client := schemas.Clients{}
	if err := config.DB.First(&client, clientID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cliente não encontrado",
		})
	}

	var total = 0.0
	allproducts := []schemas.Products{}
	for _, prod := range req.Products {

		product, err := getProductsValue(prod.ProductID, client.StoreID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if product.StockQuantity < prod.Quantity {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Quantidade de produto insuficiente",
			})
		}
		totalValues := product.Price * float64(prod.Quantity)
		total += totalValues

		// Armazenar produto com quantidade do pedido (para usar depois)
		productCopy := *product
		productCopy.StockQuantity = prod.Quantity // Quantidade sendo comprada
		allproducts = append(allproducts, productCopy)
	}
	// Usar transação para garantir consistência
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	orderNumber := generateOrderNumber()

	order := schemas.Orders{
		ClientID:      client.ID,
		OrderNumber:   orderNumber,
		Value:         total,
		Originalvalue: total,
		StatusPayment: schemas.PendingPayment,
		DeliveryType:  schemas.DeliveryType(req.DeliveryType),
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao criar pedido",
			"details": err.Error(),
		})
	}

	for _, prod := range allproducts {
		if err := createOrderItems(tx, order.ID, prod.ID, prod.StockQuantity, prod.Price, prod.Price); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Erro ao criar item do pedido",
				"details": err.Error(),
			})
		}

		// Atualizar estoque do produto
		if err := updateProductStock(tx, prod.ID, prod.StockQuantity); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Erro ao atualizar estoque do produto",
				"details": err.Error(),
			})
		}
	}

	// Commit da transação
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao finalizar pedido",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pedido criado com sucesso",
	})
}

func getProductsValue(prodID, storeID uint64) (*schemas.Products, error) {
	productSchemas := &schemas.Products{}
	product := config.DB.Where("id = ? AND store_id = ?", prodID, storeID).First(productSchemas)
	if product.Error != nil {
		return nil, fmt.Errorf("produto não encontrado")
	}
	return productSchemas, nil
}

func createOrderItems(tx *gorm.DB, orderID, productID uint64, quantity int, price, originalPrice float64) error {
	orderItem := schemas.OrderItems{
		OrderID:       orderID,
		ProductID:     productID,
		Quantity:      quantity,
		Price:         price,
		OriginalPrice: originalPrice,
	}
	if err := tx.Create(&orderItem).Error; err != nil {
		return fmt.Errorf("erro ao criar item do pedido: %w", err)
	}
	return nil
}

func updateProductStock(tx *gorm.DB, productID uint64, quantitySold int) error {
	if err := tx.Model(&schemas.Products{}).Where("id = ?", productID).
		Update("stock_quantity", gorm.Expr("stock_quantity - ?", quantitySold)).Error; err != nil {
		return fmt.Errorf("erro ao atualizar estoque: %w", err)
	}
	return nil
}

func generateOrderNumber() string {
	// Formato: ORD + timestamp + 3 dígitos aleatórios
	timestamp := time.Now().Format("20060102150405")
	randomNum := rand.Intn(900) + 100 // 100-999
	return fmt.Sprintf("ORD%s%d", timestamp, randomNum)
}
