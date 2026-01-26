package router

import (
	authcontroller "backend_camisaria_store/controller/auth"
	controller "backend_camisaria_store/controller/products"
	userController "backend_camisaria_store/controller/user"

	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App) {
	// Rotas públicas - não precisam de autenticação
	auth := app.Group("/api/auth")
	auth.Post("/login", authcontroller.LoginUser)
	auth.Post("/register", userController.CreateUser)

	// Rotas protegidas - requerem autenticação usuario
	admin := app.Group("/api/admin", authcontroller.AdminMiddlware)
	protected := app.Group("/api", authcontroller.UserMiddleware)

	// Rotas de usuários protegidas
	usersProtected := protected.Group("/users")
	usersProtected.Get("/", userController.GetUsers)      // Apenas usuários autenticados
	usersProtected.Get("/:id", userController.GetUser)    // Apenas usuários autenticados
	usersProtected.Put("/:id", userController.UpdateUser) // Apenas usuários autenticados
	admin.Delete("/:id", userController.DeleteUser)       // Apenas usuários autenticados

	// Rotas de produtos
	products := protected.Group("/products")
	products.Post("/", controller.CreateProduct)      // Criar produto
	products.Get("/", controller.ListProducts)        // Listar produtos
	products.Get("/:id", controller.GetProduct)       // Buscar produto por ID
	products.Put("/:id", controller.UpdateProduct)    // Atualizar produto
	products.Delete("/:id", controller.DeleteProduct) // Deletar produto

}
