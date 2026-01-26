package minio

import (
	"backend_camisaria_store/common"
	"backend_camisaria_store/config"
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func UploadProductImage(file *multipart.FileHeader, productName, dir string, productID uint64) (string, string, error) {
	src, err := file.Open()

	if err != nil {
		return "", "", err
	}

	defer src.Close()

	if dir != "" && !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if ext == "" {
		ext = ".png"
	}

	clenproduct := slugify(productName)

	objectName := dir + clenproduct + "-" + strconv.FormatUint(productID, 10) + "-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ext

	ct := file.Header.Get("content-Type")
	if ct == "" {
		ct = mime.TypeByExtension(ext)
	}

	_, err = config.MinioClient.PutObject(
		context.Background(),
		config.BunkedName,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: ct},
	)

	if err != nil {
		return "", "", err
	}

	publicURL := common.PublicObjectURL(objectName)
	return publicURL, objectName, nil
}

func ValidateImageFile(file *multipart.FileHeader) error {
	// Tipos MIME permitidos
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}

	// Verificar tipo MIME
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return fiber.NewError(fiber.StatusBadRequest, "Tipo de arquivo não permitido. Use apenas JPEG, PNG, WebP ou GIF")
	}

	// Verificar extensão do arquivo
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}

	if !allowedExts[ext] {
		return fiber.NewError(fiber.StatusBadRequest, "Extensão de arquivo não permitida")
	}

	// Verificar tamanho máximo (5MB)
	const maxSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxSize {
		return fiber.NewError(fiber.StatusBadRequest, "Arquivo muito grande. Tamanho máximo: 5MB")
	}

	// Verificar tamanho mínimo (1KB)
	const minSize = 1024 // 1KB
	if file.Size < minSize {
		return fiber.NewError(fiber.StatusBadRequest, "Arquivo muito pequeno. Tamanho mínimo: 1KB")
	}

	return nil
}

func ObjectKeyFormUrl(u string) (string, error) {
	if u == "" {
		return "", fmt.Errorf("campo vazio")
	}

	parsed, err := url.ParseRequestURI(u)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return strings.TrimLeft(u, "/"), nil
	}

	path := strings.TrimLeft(parsed.Path, "/")

	if b := strings.TrimSpace(config.BunkedName); b != "" {
		path = strings.TrimPrefix(path, b+"/")
	}

	key, dacerr := url.PathUnescape(path)

	if dacerr != nil {
		key = path
	}

	if key == "nil" {
		return "", fmt.Errorf("não foi possivel extrair a key do objeto")
	}

	return key, nil

}
