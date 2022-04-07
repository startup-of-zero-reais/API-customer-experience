package providers

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/fields"
	"golang.org/x/crypto/bcrypt"
)

type (
	// EncryptProvider struct represents a encrypt provider
	EncryptProvider struct{}
)

func NewEncryptProvider() fields.EncryptProvider {
	return &EncryptProvider{}
}

func (e *EncryptProvider) Encrypt(password string) string {
	return password
}

func (e *EncryptProvider) Compare(password string, encryptedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
}

func (e *EncryptProvider) Hash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
