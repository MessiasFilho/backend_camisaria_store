package controller

import (
	"errors"
	"strings"
)

type DeliveryType string

const (
	PickupDelivery   DeliveryType = "pickup"
	DeliveryDelivery DeliveryType = "delivery"
)

type CreateOrderRequest struct {
	Products     []productsStruct `json:"products" validate:"required,min=1"`
	DeliveryType DeliveryType     `json:"delivery_type" validate:"required,oneof=pickup delivery"`
}

type productsStruct struct {
	ProductID uint64 `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

func (req *CreateOrderRequest) Validate() error {
	var errs []string

	// Validação dos Produtos
	if len(req.Products) == 0 {
		errs = append(errs, "produtos são obrigatórios")
	}

	for _, product := range req.Products {
		if product.ProductID == 0 {
			errs = append(errs, "product_id é obrigatório")
		}
		if product.Quantity <= 0 {
			errs = append(errs, "quantity deve ser maior que zero")
		}
	}

	// Validação do DeliveryType
	if req.DeliveryType == "" {
		errs = append(errs, "delivery_type é obrigatório")
	}

	// Se houver erros, retorna um erro composto
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
