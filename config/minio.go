package config

import (
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioClient *minio.Client
)

var (
	BunkedName string
)

func InitMinio() error {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ROOT_USER")
	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	BunkedName = os.Getenv("MINIO_BUCKET")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("erro ao conectar no minio: %v", err)
	}
	return nil

}
