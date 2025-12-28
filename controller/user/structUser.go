package controller

import (
	"errors"
	"regexp"
	"strings"
)

// CreateUserRequest representa a requisição para criar um usuário
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest representa a requisição para atualizar um usuário
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
}

// UserResponse representa a resposta da API para usuários
type UserResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UserListResponse representa a resposta paginada para listagem de usuários
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Pages int            `json:"pages"`
}

// Validate realiza validação personalizada dos campos do CreateUserRequest
func (req *CreateUserRequest) Validate() error {
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
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.Email) {
			errs = append(errs, "e-mail deve ter um formato válido")
		}
	}

	// Validação da Senha
	if strings.TrimSpace(req.Password) == "" {
		errs = append(errs, "senha é obrigatória")
	} else if len(req.Password) < 6 {
		errs = append(errs, "senha deve ter pelo menos 6 caracteres")
	}

	// Se houver erros, retorna um erro composto
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// Validate realiza validação personalizada dos campos do UpdateUserRequest
func (req *UpdateUserRequest) Validate() error {
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
		} else {
			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
			if !emailRegex.MatchString(*req.Email) {
				errs = append(errs, "e-mail deve ter um formato válido")
			}
		}
	}

	// Validação da Senha (se informada)
	if req.Password != nil {
		if strings.TrimSpace(*req.Password) == "" {
			errs = append(errs, "senha não pode ser vazia")
		} else if len(*req.Password) < 6 {
			errs = append(errs, "senha deve ter pelo menos 6 caracteres")
		}
	}

	// Se houver erros, retorna um erro composto
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
