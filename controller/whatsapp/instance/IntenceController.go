package controller

import (
	"backend_camisaria_store/config"
	"backend_camisaria_store/schemas"
	whatsapp "backend_camisaria_store/service/whatsapp/config"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateInstance(c *fiber.Ctx) error {
	req := CreateInstanceRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao processar dados da requisição",
		})
	}
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var existing schemas.Instance
	switch err := config.DB.Where("hash = ?", &existing.Hash).First(&existing).Error; {
	case err == nil:
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Já existe uma instância com este nome",
		})
	case !errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao verificar instâncias cadastradas",
		})
	}

	baseURL, apiKey := whatsapp.GetBaseURLAndAPIKey()
	if baseURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "WHATSAPP_BASE_URL não configurada",
		})
	}
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "WHATSAPP_API_KEY não configurada",
		})
	}

	createURL := baseURL + "/instance/create"

	payload := map[string]interface{}{
		"instanceName": req.Name,
		"integration":  req.Integration,
		"qrcode":       req.Qrcode,
	}
	bodyJSON, err := json.Marshal(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao montar JSON",
		})
	}

	httpReq, err := http.NewRequest(http.MethodPost, createURL, bytes.NewReader(bodyJSON))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	httpReq.Header.Set("Content-Type", "application/json")
	whatsapp.SetAuthHeaders(httpReq, apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Falha ao contatar serviço WhatsApp",
			"details": err.Error(),
			"url":     createURL,
		})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Erro ao ler resposta do WhatsApp",
		})
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error":   "Erro ao criar instância na Evolution API",
			"details": string(respBody),
			"url":     createURL,
		})
	}

	hash, err := whatsapp.UniqueHashName()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	instance := schemas.Instance{
		Name:        req.Name,
		Integration: req.Integration,
		QrCode:      req.Qrcode,
		Status:      "active",
		Hash:        hash,
	}
	if err := config.DB.Create(&instance).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao salvar instância",
		})
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "" {
		c.Set("Content-Type", ct)
	} else {
		c.Type("json")
	}
	return c.Status(resp.StatusCode).Send(respBody)
}
