package router

import (
	authcontroller "backend_camisaria_store/controller/auth"
	userController "backend_camisaria_store/controller/user"

	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App) {
	// Rotas públicas - não precisam de autenticação
	auth := app.Group("/api/auth")
	auth.Post("/login", authcontroller.LoginUser)

	// Rotas de usuários - algumas podem ser públicas, outras protegidas
	userPath := "/api/users"
	user := app.Group(userPath)
	user.Post("/", userController.CreateUser) // Criar usuário (público)

	// Rotas protegidas - requerem autenticação
	protected := app.Group("/api", authcontroller.AuthMiddleware)

	// Rotas de usuários protegidas
	usersProtected := protected.Group("/users")
	usersProtected.Get("/", userController.GetUsers)         // Apenas usuários autenticados
	usersProtected.Get("/:id", userController.GetUser)       // Apenas usuários autenticados
	usersProtected.Put("/:id", userController.UpdateUser)    // Apenas usuários autenticados
	usersProtected.Delete("/:id", userController.DeleteUser) // Apenas usuários autenticados

	// Rotas administrativas - apenas para admins
	admin := protected.Group("/admin", authcontroller.RequireAdmin)
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Admin Dashboard",
			"user": fiber.Map{
				"id": c.Locals("userID"),
				"name": c.Locals("userName"),
				"role": c.Locals("userRole"),
			},
		})
	})

	// Rotas para clientes - apenas para usuários com role "client"
	client := protected.Group("/client", authcontroller.RequireRole("client"))
	client.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Client Profile",
			"user": fiber.Map{
				"id": c.Locals("userID"),
				"name": c.Locals("userName"),
				"email": c.Locals("userEmail"),
				"type": c.Locals("userType"),
			},
		})
	})

	// Rota de exemplo para usuários normais
	userRoutes := protected.Group("/user", authcontroller.RequireRole("user"))
	userRoutes.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "User Profile",
			"user": fiber.Map{
				"id": c.Locals("userID"),
				"name": c.Locals("userName"),
				"type": c.Locals("userType"),
			},
		})
	})
}
