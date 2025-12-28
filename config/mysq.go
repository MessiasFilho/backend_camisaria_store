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

	err := godotenv.Load()

	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	urlDB := os.Getenv("DB")

	// Configuração padrão para desenvolvimento local
	if urlDB == "" {
		log.Println("Variável DB não encontrada, usando configuração padrão para desenvolvimento")
		urlDB = "root:@tcp(localhost:3306)/loja_camisaria?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(urlDB), &gorm.Config{})

	if err != nil {
		fmt.Println("Erro ao conectar ao MySQL:", err) // Corrigido para fmt.Println
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
