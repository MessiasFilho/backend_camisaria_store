package router

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Initialize() {
	app := fiber.New(fiber.Config{
		AppName: "Santiago store backend",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Authorization,Accept",
	}))

	app.Use(func(c *fiber.Ctx) error {
		// Obter o esquema (http ou https)
		scheme := c.Protocol()

		// Obter o hostname
		host := c.Hostname()

		// Obter o caminho original da URL
		path := c.OriginalURL()

		// Concatenar para formar a URL completa
		fullURL := fmt.Sprintf("%s://%s%s", scheme, host, path)

		// Obter o IP que está solicitando a request
		ip := c.IP()

		// Obter o User-Agent do solicitante
		userAgent := c.Get("User-Agent")

		// Imprimir as informações no console
		fmt.Printf("URL: %s, IP: %s, User-Agent: %s\n", fullURL, ip, userAgent)

		// Continue para a próxima rota
		return c.Next()
	})

	InitializeRoutes(app)

	app.Listen(":4041")
}
