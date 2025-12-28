package controller

import (
	"backend_camisaria_store/common"
	"backend_camisaria_store/config"
	"backend_camisaria_store/schemas"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LoginUser(c *fiber.Ctx) error {
	data := LoginStruct{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao realizar login",
		})
	}

	if err := data.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user = schemas.Users{}
	db := config.DB

	result := db.Where(&schemas.Users{
		Email: data.Email,
	}).First(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email ou senha inválidos.",
		})
	}

	if !common.CheckPassword(data.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email ou senha inválidos",
		})
	}

	line := 2

	if user.Role == "admin" {
		line = 1
	}

	claims := jwt.MapClaims{
		"code": user.ID,
		"name": user.Name,
		"line": line,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"acees_token": t})
}
