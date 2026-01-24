package minio

import (
	"backend_camisaria_store/common"
	"backend_camisaria_store/config"
	"context"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
