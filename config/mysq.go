package config

import (
	"backend_camisaria_store/schemas"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeMysql() (*gorm.DB, error) {

	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	host := mustGetEnv("DB_HOST")
	port := mustGetEnv("DB_PORT")
	user := mustGetEnv("DB_USER")
	pass := mustGetEnv("DB_PASSWORD")
	name := mustGetEnv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if os.Getenv("AUTO_MIGRATE") == "true" {
		if err := db.AutoMigrate(
			&schemas.Users{},
			&schemas.Clients{},
			&schemas.Orders{},
			&schemas.Products{},
			&schemas.ProductImage{},
		); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Variável de ambiente obrigatória não definida: %s", key)
	}
	return value
}
