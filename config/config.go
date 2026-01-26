package config

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Init() error {
	var err error

	// Inicializar MySQL
	DB, err = InitializeMysql()
	if err != nil {
		return fmt.Errorf("error initialize mysql %v", err)
	}

	// Inicializar MinIO
	err = InitMinio()
	if err != nil {
		return fmt.Errorf("error initialize minio %v", err)
	}

	return nil
}
