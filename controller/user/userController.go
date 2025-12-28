package controller

import (
	"backend_camisaria_store/config"
	"backend_camisaria_store/schemas"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	req := CreateUserRequest{}

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

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao processar senha",
			"details": "Falha ao criar hash da senha",
		})
	}

	// Mapear request para schema
	user := schemas.Users{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Salvar no banco
	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao salvar usuário",
			"details": err.Error(),
		})
	}

	// Retornar resposta de sucesso
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created",
	})
}

func GetUsers(c *fiber.Ctx) error {
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

	var users []schemas.Users
	var total int64

	// Buscar total de usuários
	if err := config.DB.Model(&schemas.Users{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao contar usuários",
			"details": err.Error(),
		})
	}

	// Buscar usuários com paginação
	if err := config.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao buscar usuários",
			"details": err.Error(),
		})
	}

	// Calcular total de páginas
	pages := int((total + int64(limit) - 1) / int64(limit))

	// Converter para response
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return c.JSON(UserListResponse{
		Users: userResponses,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pages,
	})
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Converter string para uint
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	var user schemas.Users
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuário não encontrado",
		})
	}

	return c.JSON(UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Converter string para uint
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	req := UpdateUserRequest{}

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

	// Verificar se usuário existe
	var user schemas.Users
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuário não encontrado",
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

	// Executar update
	if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao atualizar usuário",
			"details": err.Error(),
		})
	}

	// Buscar usuário atualizado
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao buscar usuário atualizado",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Usuário atualizado com sucesso",
		"user": UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Converter string para uint
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	// Verificar se usuário existe
	var user schemas.Users
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuário não encontrado",
		})
	}

	// Deletar usuário
	if err := config.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao deletar usuário",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Usuário deletado com sucesso",
	})
}
