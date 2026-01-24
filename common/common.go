package common

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func PublicObjectURL(objectName string) string {
	scheme := "http"
	if os.Getenv("MINIO_USE_SSL") == "true" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s",
		scheme,
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_BUCKET"),
		objectName,
	)
}

func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("your-secret-key-change-in-production")
	}
	return []byte(secret)
}
