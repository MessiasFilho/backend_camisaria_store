package controller

import (
	"backend_camisaria_store/common"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware valida tokens JWT e gerencia diferentes tipos de usuários (admin, user, client)
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token de autenticação não fornecido",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Formato de token inválido. Use: Bearer <token>",
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
	userName, _ := claims["name"].(string)
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

	// Armazenar informações do usuário no contexto para uso nas rotas
	c.Locals("userID", uint(userIDFloat))
	c.Locals("userType", userType)
	c.Locals("userRole", userRole)
	c.Locals("userName", userName)
	c.Locals("userEmail", userEmail)

	// Log opcional para debug (remover em produção)
	// fmt.Printf("Usuário autenticado: ID=%d, Type=%s, Role=%s\n", uint(userIDFloat), userType, userRole)

	return c.Next()
}

// RequireRole é um middleware auxiliar para verificar roles específicos
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Usuário não autenticado",
			})
		}

		role, ok := userRole.(string)
		if !ok || role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Acesso negado. Role necessário: " + requiredRole,
			})
		}

		return c.Next()
	}
}

// RequireAdmin é um atalho para RequireRole("admin")
func RequireAdmin(c *fiber.Ctx) error {
	return RequireRole("admin")(c)
}
