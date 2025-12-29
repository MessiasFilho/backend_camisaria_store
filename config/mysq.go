package config

import (
	"backend_camisaria_store/schemas"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeMysql() (*gorm.DB, error) {

	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
		port = "3306"
		user = "root"
		name = "loja_camisaria"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&schemas.Users{},
		&schemas.Clients{},
		&schemas.Orders{},
		&schemas.Products{},
		&schemas.ProductImage{},
	)

	if err != nil {
		fmt.Printf("Mysql automigration error %v", err)
	}

	return db, err

}
