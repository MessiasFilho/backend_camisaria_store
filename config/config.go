package config

import (
	"backend_camisaria_store/schemas"
	whatsapp "backend_camisaria_store/service/whatsapp/config"
	"fmt"
	"log"

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

	// Inicializar instância default
	err = InitInstanceDbDefault()
	if err != nil {
		return fmt.Errorf("error initialize instance default %v", err)
	}

	return nil
}

// InitInstanceDbDefault cria a instância "default" somente quando a tabela instances está vazia.
func InitInstanceDbDefault() error {
	var count int64
	if err := DB.Model(&schemas.Instance{}).Count(&count).Error; err != nil {
		return fmt.Errorf("error count instances: %w", err)
	}
	if count > 0 {
		return nil
	}

	instance := schemas.Instance{
		Name:        "default",
		Integration: "WHATSAPP-BAILEYS",
		QrCode:      true,
		Status:      "active",
		Hash:        whatsapp.GenerateHashInstance(),
	}
	if err := DB.Create(&instance).Error; err != nil {
		return fmt.Errorf("error create instance default: %w", err)
	}

	log.Println("instância WhatsApp padrão criada (name=default)")
	return nil
}
