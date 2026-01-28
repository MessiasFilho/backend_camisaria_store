package controller

import (
	"backend_camisaria_store/config"
	"backend_camisaria_store/schemas"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateClient(c *fiber.Ctx) error {
	req := CreateClientRequest{}

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

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao processar senha",
		})
	}

	// Mapear request para schema
	client := schemas.Clients{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
	}

	// Mapear endereço apenas se fornecido
	if req.Address != nil {
		client.Address = schemas.Address{
			Street:       req.Address.Street,
			Number:       req.Address.Number,
			Complement:   req.Address.Complement,
			Neighborhood: req.Address.Neighborhood,
			City:         req.Address.City,
			State:        req.Address.State,
			ZipCode:      req.Address.ZipCode,
		}
	}

	// Salvar no banco
	if err := config.DB.Create(&client).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao salvar cliente",
		})
	}

	// Retornar resposta de sucesso
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Cliente criado com sucesso",
	})

}

func GetClients(c *fiber.Ctx) error {
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

	var clients []schemas.Clients
	var total int64

	// Buscar total de clientes
	if err := config.DB.Model(&schemas.Clients{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao contar clientes",
			"details": err.Error(),
		})
	}

	// Buscar clientes com paginação
	if err := config.DB.Offset(offset).Limit(limit).Find(&clients).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao buscar clientes",
			"details": err.Error(),
		})
	}

	// Calcular total de páginas
	pages := int((total + int64(limit) - 1) / int64(limit))

	// Converter para response
	var clientResponses []ClientResponse
	for _, client := range clients {
		var address *Address
		if client.Address.Street != "" {
			address = &Address{
				Street:       client.Address.Street,
				Number:       client.Address.Number,
				Complement:   client.Address.Complement,
				Neighborhood: client.Address.Neighborhood,
				City:         client.Address.City,
				State:        client.Address.State,
				ZipCode:      client.Address.ZipCode,
			}
		}

		clientResponses = append(clientResponses, ClientResponse{
			ID:        client.ID,
			Name:      client.Name,
			Email:     client.Email,
			Phone:     client.Phone,
			Status:    client.Status,
			Role:      string(client.Role),
			Address:   address,
			CreatedAt: client.CreatedAt,
			UpdatedAt: client.UpdatedAt,
		})
	}

	return c.JSON(ClientListResponse{
		Clients: clientResponses,
		Total:   total,
		Page:    page,
		Limit:   limit,
		Pages:   pages,
	})
}

func GetClient(c *fiber.Ctx) error {
	id := c.Params("id")

	// Converter string para uint
	clientID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	var client schemas.Clients
	if err := config.DB.First(&client, clientID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cliente não encontrado",
		})
	}

	var address *Address
	if client.Address.Street != "" {
		address = &Address{
			Street:       client.Address.Street,
			Number:       client.Address.Number,
			Complement:   client.Address.Complement,
			Neighborhood: client.Address.Neighborhood,
			City:         client.Address.City,
			State:        client.Address.State,
			ZipCode:      client.Address.ZipCode,
		}
	}

	return c.JSON(ClientResponse{
		ID:        client.ID,
		Name:      client.Name,
		Email:     client.Email,
		Phone:     client.Phone,
		Status:    client.Status,
		Role:      string(client.Role),
		Address:   address,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	})
}

func UpdateClient(c *fiber.Ctx) error {
	id := c.Params("id")

	// Converter string para uint
	clientID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	req := UpdateClientRequest{}
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

	// Verificar se cliente existe
	var client schemas.Clients
	if err := config.DB.First(&client, clientID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cliente não encontrado",
		})
	}

	// Atualizar campos se fornecidos
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Password != nil {
		// Hash da nova senha
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Erro ao processar senha",
				"details": "Falha ao criar hash da senha",
			})
		}
		updates["password"] = string(hashedPassword)
	}

	// Atualizar endereço se fornecido
	if req.Address != nil {
		updates["address"] = schemas.Address{
			Street:       req.Address.Street,
			Number:       req.Address.Number,
			Complement:   req.Address.Complement,
			Neighborhood: req.Address.Neighborhood,
			City:         req.Address.City,
			State:        req.Address.State,
			ZipCode:      req.Address.ZipCode,
		}
	}

	// Executar update
	if err := config.DB.Model(&client).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao atualizar cliente",
			"details": err.Error(),
		})
	}

	// Buscar cliente atualizado
	if err := config.DB.First(&client, clientID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao buscar cliente atualizado",
			"details": err.Error(),
		})
	}

	var address *Address
	if client.Address.Street != "" {
		address = &Address{
			Street:       client.Address.Street,
			Number:       client.Address.Number,
			Complement:   client.Address.Complement,
			Neighborhood: client.Address.Neighborhood,
			City:         client.Address.City,
			State:        client.Address.State,
			ZipCode:      client.Address.ZipCode,
		}
	}

	return c.JSON(fiber.Map{
		"message": "Cliente atualizado com sucesso",
		"client": ClientResponse{
			ID:        client.ID,
			Name:      client.Name,
			Email:     client.Email,
			Phone:     client.Phone,
			Status:    client.Status,
			Role:      string(client.Role),
			Address:   address,
			CreatedAt: client.CreatedAt,
			UpdatedAt: client.UpdatedAt,
		},
	})
}

func DeleteClient(c *fiber.Ctx) error {
	id := c.Params("id")

	// Converter string para uint
	clientID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	// Verificar se cliente existe
	var client schemas.Clients
	if err := config.DB.First(&client, clientID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cliente não encontrado",
		})
	}

	// Soft delete (desativar ao invés de deletar)
	if err := config.DB.Model(&client).Update("status", "inactive").Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao deletar cliente",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cliente deletado com sucesso",
	})
}
