package controller

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CreateInstanceRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Integration string `json:"integration" validate:"required,min=3,max=100"`
	Qrcode      bool   `json:"qrcode"`
}

func (req *CreateInstanceRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return err
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New("nome é obrigatório")
	}
	if strings.TrimSpace(req.Integration) == "" {
		return errors.New("integração é obrigatória")
	}
	return nil
}
