package config

import (
	"backend_camisaria_store/schemas"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeMysql() (*gorm.DB, error) {

	// Carrega .env apenas fora de produção
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	env := getEnv("ENV", "development")

	host := mustGetEnv("DB_HOST")
	port := mustGetEnv("DB_PORT")
	user := mustGetEnv("DB_USER")
	name := mustGetEnv("DB_NAME")

	// 🔐 Regra especial para senha
	pass := os.Getenv("DB_PASSWORD")

	if pass == "" && env != "development" {
		log.Fatal("Variável de ambiente obrigatória não definida: DB_PASSWORD")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		pass, // pode ser vazio apenas em DEV
		host,
		port,
		name,
	)

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}

	// Auto migrate controlado por env

	if err := db.AutoMigrate(
		&schemas.Users{},
		&schemas.Clients{},
		&schemas.Orders{},
		&schemas.Products{},
		&schemas.OrderItems{},
		&schemas.Address{},
		&schemas.Instance{},
	); err != nil {
		return nil, err
	}

	return db, nil
}

// ========================
// Helpers
// ========================

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Variável de ambiente obrigatória não definida: %s", key)
	}
	return value
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
