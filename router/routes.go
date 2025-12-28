package router

import (
	authcontroller "backend_camisaria_store/controller/auth"
	userController "backend_camisaria_store/controller/user"

	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App) {
	// Rotas de Usuários
	userPath := "/api/users"
	authPath := "/api/auth"
	user := app.Group(userPath)
	auth := app.Group(authPath)
	user.Post("/", userController.CreateUser)      // Criar usuário
	user.Get("/", userController.GetUsers)         // Listar usuários com paginação
	user.Get("/:id", userController.GetUser)       // Buscar usuário por ID
	user.Put("/:id", userController.UpdateUser)    // Atualizar usuário
	user.Delete("/:id", userController.DeleteUser) // Deletar usuário

	auth.Post("/", authcontroller.LoginUser)
}
