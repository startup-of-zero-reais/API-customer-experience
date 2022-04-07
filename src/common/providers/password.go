package providers

import (
	"golang.org/x/crypto/bcrypt"
)

type (
	// EncryptProvider is a interface to encrypt and decrypt data
	EncryptProvider interface {
		Hash(password string) string
		Compare(password string, hashedPassword string) error
	}

	// EncryptProviderImpl struct represents a encrypt provider
	EncryptProviderImpl struct{}
)

func NewEncryptProvider() EncryptProvider {
	return &EncryptProviderImpl{}
}

func (e *EncryptProviderImpl) Encrypt(password string) string {
	return password
}

func (e *EncryptProviderImpl) Compare(password string, encryptedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
}

func (e *EncryptProviderImpl) Hash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
