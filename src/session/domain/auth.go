package domain

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
)

type (
	Password interface {
		Hash() string
		Validate(password string) error
	}

	PasswordImpl struct {
		hash            string
		encryptProvider providers.EncryptProvider
	}

	// Auth struct represents a auth entity
	Auth struct {
		Token    string   `json:"token"`
		Password Password `json:"password"`
	}

	AuthInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

func NewPassword(hash string) Password {
	return &PasswordImpl{
		hash:            hash,
		encryptProvider: providers.NewEncryptProvider(),
	}
}

func (p *PasswordImpl) Hash() string {
	return p.hash
}

func (p *PasswordImpl) Validate(password string) error {
	err := p.encryptProvider.Compare(password, p.hash)
	if err != nil {
		return validation.UnauthorizedError("credenciais inválidas")
	}

	return nil
}
