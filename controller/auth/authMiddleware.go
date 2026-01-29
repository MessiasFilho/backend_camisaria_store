package controller

import (
	"backend_camisaria_store/common"
	"backend_camisaria_store/config"
	"backend_camisaria_store/schemas"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AuthMiddleware valida tokens JWT e verifica usuário no banco de dados
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Token de autenticação não fornecido",
			"message": "Inclua o header 'Authorization: Bearer <token>' na requisição",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Formato de token inválido",
			"message": "O token deve começar com 'Bearer '",
		})
	}

	// Parse e valida o token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Método de assinatura inválido")
		}
		return common.GetJWTSecret(), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token inválido: " + err.Error(),
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token expirado ou inválido",
		})
	}

	// Extrair claims do token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Estrutura do token inválida",
		})
	}

	// Extrair informações do usuário dos claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "ID do usuário não encontrado no token",
		})
	}

	userType, _ := claims["user_type"].(string)
	userRole, _ := claims["role"].(string)
	userEmail, _ := claims["email"].(string)

	// Se user_type estiver vazio, usar role como fallback
	if userType == "" {
		userType = userRole
	}

	// Validar que o role é válido
	switch userType {
	case "admin", "user", "client":
		// Role válido, continuar
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Tipo de usuário inválido",
		})
	}

	// ====================
	// Verificar se o usuário ainda existe e está ativo no banco de dados
	// ====================

	userID := uint(userIDFloat)
	var dbUser schemas.Users

	result := config.DB.Where("id = ?", userID).First(&dbUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Usuário não encontrado",
				"message": "O usuário associado a este token não existe mais",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro interno do servidor",
			"message": "Erro ao validar usuário",
		})
	}

	// Verificar se as informações do token batem com o banco
	if dbUser.Email != userEmail {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Token inválido",
			"message": "As informações do token não correspondem aos dados do usuário",
		})
	}

	if string(dbUser.Role) != userType {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Token inválido",
			"message": "As permissões do usuário foram alteradas",
		})
	}

	// Usar dados atualizados do banco (não do token)
	c.Locals("user_id", dbUser.ID)
	c.Locals("user_type", string(dbUser.Role))
	c.Locals("user_role", string(dbUser.Role))
	c.Locals("user_name", dbUser.Name)
	c.Locals("user_email", dbUser.Email)

	// Usuário autenticado com sucesso

	return c.Next()
}

func AdminMiddlware(c *fiber.Ctx) error {
	// Verificar se o usuário está autenticado
	userTypeLocal := c.Locals("user_type")
	if userTypeLocal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Autenticação necessária",
			"message": "Você precisa estar logado para acessar este recurso",
		})
	}

	// Fazer type assertion segura
	userType, ok := userTypeLocal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro interno",
			"message": "Erro ao processar dados de autenticação",
		})
	}

	// Verificar se é admin
	if userType != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Acesso negado",
			"message": "Apenas administradores podem acessar este recurso",
		})
	}

	return c.Next()
}

func UserMiddleware(c *fiber.Ctx) error {
	// Verificar se o usuário está autenticado (se user_type existe)
	userTypeLocal := c.Locals("user_type")
	if userTypeLocal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Autenticação necessária",
			"message": "Você precisa estar logado para acessar este recurso",
		})
	}

	// Fazer type assertion segura
	userType, ok := userTypeLocal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro interno",
			"message": "Erro ao processar dados de autenticação",
		})
	}

	// Verificar permissões
	if userType != "user" && userType != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Acesso negado",
			"message": "Apenas usuários e administradores podem acessar este recurso",
		})
	}

	return c.Next()
}

func ClientMiddleware(c *fiber.Ctx) error {
	// Verificar se o usuário está autenticado
	userTypeLocal := c.Locals("user_type")
	if userTypeLocal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Autenticação necessária",
			"message": "Você precisa estar logado para acessar este recurso",
		})
	}

	// Fazer type assertion segura
	userType, ok := userTypeLocal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro interno",
			"message": "Erro ao processar dados de autenticação",
		})
	}
	
	// Verificar se é cliente
	if userType != "client" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Acesso negado",
			"message": "Apenas clientes podem acessar este recurso",
		})
	}

	return c.Next()
}
