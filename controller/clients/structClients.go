package controller

import (
	"errors"
	"strings"
	"time"
)

type CreateClientRequest struct {
	Name     string   `json:"name" validate:"required,min=3,max=255"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
	Phone    string   `json:"phone" validate:"required,min=3,max=255"`
	Address  *Address `json:"address,omitempty"`
}

type Address struct {
	Street       string `json:"street" validate:"required,min=3,max=255"`
	Number       string `json:"number" validate:"required,min=1,max=10"`
	Complement   string `json:"complement" validate:"required,min=3,max=255"`
	Neighborhood string `json:"neighborhood" validate:"required,min=3,max=255"`
	City         string `json:"city" validate:"required,min=3,max=255"`
	State        string `json:"state" validate:"required,min=3,max=255"`
	ZipCode      string `json:"zip_code" validate:"required,min=3,max=255"`
}

// UpdateClientRequest representa a requisição para atualizar um cliente
type UpdateClientRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Email    *string  `json:"email,omitempty" validate:"omitempty,email"`
	Password *string  `json:"password,omitempty" validate:"omitempty,min=6"`
	Phone    *string  `json:"phone,omitempty" validate:"omitempty,min=3,max=255"`
	Status   *string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
	Address  *Address `json:"address,omitempty"`
}

// ClientResponse representa a resposta da API para cliente
type ClientResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
	Address   *Address  `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ClientListResponse representa a resposta paginada para listagem de clientes
type ClientListResponse struct {
	Clients []ClientResponse `json:"clients"`
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
	Pages   int              `json:"pages"`
}

func (req *CreateClientRequest) Validate() error {
	var errs []string

	// Validação do Nome
	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, "nome é obrigatório")
	} else if len(req.Name) < 3 || len(req.Name) > 255 {
		errs = append(errs, "nome deve ter entre 3 e 255 caracteres")
	}

	// Validação do Email
	if strings.TrimSpace(req.Email) == "" {
		errs = append(errs, "e-mail é obrigatório")
	} else if len(req.Email) < 3 || len(req.Email) > 255 {
		errs = append(errs, "e-mail deve ter entre 3 e 255 caracteres")
	}

	// Validação da Senha
	if strings.TrimSpace(req.Password) == "" {
		errs = append(errs, "senha é obrigatória")
	} else if len(req.Password) < 6 {
		errs = append(errs, "senha deve ter pelo menos 6 caracteres")
	}

	// Validação do Telefone
	if strings.TrimSpace(req.Phone) == "" {
		errs = append(errs, "telefone é obrigatório")
	} else if len(req.Phone) < 3 || len(req.Phone) > 255 {
		errs = append(errs, "telefone deve ter entre 3 e 255 caracteres")
	}

	// Validação do Endereço (opcional)
	// Endereço não é mais obrigatório

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil // Retorna nil se todas as validações passarem
}

func (req *UpdateClientRequest) Validate() error {
	var errs []string

	// Validação do Nome (se informado)
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			errs = append(errs, "nome não pode ser vazio")
		} else if len(*req.Name) < 3 || len(*req.Name) > 255 {
			errs = append(errs, "nome deve ter entre 3 e 255 caracteres")
		}
	}

	// Validação do Email (se informado)
	if req.Email != nil {
		if strings.TrimSpace(*req.Email) == "" {
			errs = append(errs, "e-mail não pode ser vazio")
		} else if len(*req.Email) < 3 || len(*req.Email) > 255 {
			errs = append(errs, "e-mail deve ter entre 3 e 255 caracteres")
		}
	}

	// Validação da Senha (se informada)
	if req.Password != nil && len(*req.Password) < 6 {
		errs = append(errs, "senha deve ter pelo menos 6 caracteres")
	}

	// Validação do Telefone (se informado)
	if req.Phone != nil {
		if strings.TrimSpace(*req.Phone) == "" {
			errs = append(errs, "telefone não pode ser vazio")
		} else if len(*req.Phone) < 3 || len(*req.Phone) > 255 {
			errs = append(errs, "telefone deve ter entre 3 e 255 caracteres")
		}
	}

	// Validação do Status (se informado)
	if req.Status != nil && *req.Status != "active" && *req.Status != "inactive" {
		errs = append(errs, "status deve ser 'active' ou 'inactive'")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil // Retorna nil se todas as validações passarem
}
