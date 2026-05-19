package router

import (
	authcontroller "backend_camisaria_store/controller/auth"
	clientController "backend_camisaria_store/controller/clients"
	controller "backend_camisaria_store/controller/products"
	userController "backend_camisaria_store/controller/user"
	instanceController "backend_camisaria_store/controller/whatsapp/instance"
	"backend_camisaria_store/service/minio"

	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App) {

	// para evitar conflito de middlewares

	// Rotas públicas - não precisam de autenticação

	public := app.Group("/public")
	public.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "API is running",
		})
	})

	public.Post("/login", authcontroller.LoginUser)
	public.Post("/register", userController.CreateUser)

	// Loja pública — produtos publicados na página principal
	public.Get("/store/products", controller.ListPublishedProducts)

	// Rotas protegidas - requerem autenticação
	admin := app.Group("/api/admin", authcontroller.AuthMiddleware, authcontroller.AdminMiddlware)
	admin.Post("/users", userController.CreateStaffUser) // Criar admin/user interno (Postman)
	admin.Delete("/:id", userController.DeleteUser)      // Apenas admins podem deletar

	// Grupo geral para /api/* (exceto /api/auth/* que já foi definido acima)
	protected := app.Group("/api", authcontroller.AuthMiddleware, authcontroller.UserMiddleware)

	// Rotas de usuários protegidas
	usersProtected := protected.Group("/users")
	usersProtected.Get("/", userController.GetUsers)      // Apenas usuários autenticados
	usersProtected.Get("/:id", userController.GetUser)    // Apenas usuários autenticados
	usersProtected.Put("/:id", userController.UpdateUser) // Apenas usuários autenticados

	// Rotas de produtos (rotas fixas antes de /:id)
	products := protected.Group("/products")
	products.Get("/categories", controller.ListCategoriesSummary)
	products.Get("/category/:category", controller.ListProductsByCategory)
	products.Post("/", controller.CreateProduct)
	products.Get("/", controller.ListProducts)
	products.Get("/:id", controller.GetProduct)
	products.Put("/:id", controller.UpdateProduct)    // Atualizar produto
	products.Delete("/:id", controller.DeleteProduct) // Deletar produto

	products.Post("/:id", minio.UploadImgesProduct)

	// Rota para deletar múltiplas imagens (envia lista de URLs no body)
	products.Post("/delete-images", minio.DeleteImagesMinio)

	// Rotas de clientes
	clients := protected.Group("/clients")
	clients.Post("/", clientController.CreateClient)      // Criar cliente
	clients.Get("/", clientController.GetClients)         // Listar clientes
	clients.Get("/:id", clientController.GetClient)       // Buscar cliente por ID
	clients.Put("/:id", clientController.UpdateClient)    // Atualizar cliente
	clients.Delete("/:id", clientController.DeleteClient) // Deletar cliente (soft delete)

	// Rotas de instâncias
	instances := admin.Group("/instances")
	instances.Post("/", instanceController.CreateInstance)

}
