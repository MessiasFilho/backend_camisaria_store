package controller

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type LoginStruct struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=20"`
}

func (req *LoginStruct) Validate() error {
	// Cria uma instância do validador
	validate := validator.New()

	// Valida a struct UserRequest com base nas tags de validação especificadas em cada campo
	err := validate.Struct(req)
	if err != nil {
		// Se houver erros de validação, percorre cada erro e formata uma mensagem de erro
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("o campo %s é invalido: %s", err.Field(), err.ActualTag())
			// Retorna uma mensagem de erro com o nome do campo e a regra de validação que falhou
		}
	}

	// Registra a validação personalizada para a senha

	return nil // Retorna nil se todas as validações passarem
}
