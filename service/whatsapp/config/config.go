package whatsapp

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
)

func GetBaseURLAndAPIKey() (string, string) {
	return NormalizeBaseURL(os.Getenv("WHATSAPP_BASE_URL")), NormalizeAPIKey(os.Getenv("WHATSAPP_API_KEY"))
}

func NormalizeBaseURL(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.Trim(s, `"'`)
	return strings.TrimSuffix(s, "/")
}

func NormalizeAPIKey(raw string) string {
	return strings.TrimSpace(strings.Trim(raw, `"'`))
}

func SetAuthHeaders(req *http.Request, apiKey string) {
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
}

// GenerateHashInstance gera identificador único (32 hex, estilo Evolution API).
func GenerateHashInstance() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("falha ao gerar hash da instância: " + err.Error())
	}
	return strings.ToUpper(hex.EncodeToString(b))
}

// HashFromCreateResponse extrai o hash retornado pela Evolution API, se existir.
func UniqueHashName() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)[:8]), nil
}
