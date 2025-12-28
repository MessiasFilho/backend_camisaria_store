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
	DB, err = InitializeMysql()

	if err != nil {
		return fmt.Errorf("error initialize mysql %v", err)
	}
	return nil
}
