package fields

import "errors"

type (
	// EncryptProvider is a interface to encrypt and decrypt data
	EncryptProvider interface {
		Hash(password string) string
		Compare(password string, hashedPassword string) error
	}

	// Password is a struct that contains the password and the EncrypyProvider
	Password struct {
		EncryptProvider
		password string
		hash     string
	}
)

// NewPassword is the constructor of Password
func NewPassword(e EncryptProvider, password string) *Password {
	return &Password{
		EncryptProvider: e,
		password:        password,
	}
}

// IsValid is a method that checks if the password is valid
func (p *Password) IsValid() []string {
	var errs []string
	if p.password == "" {
		errs = append(errs, "o campo de senha é obrigatório")
	}

	if len(p.password) < 6 {
		errs = append(errs, "o campo de senha deve ter no mínimo 6 caracteres")
	}

	return errs
}

func (p *Password) Encrypt() *Password {
	p.hash = p.EncryptProvider.Hash(p.password)
	return p
}

// Hash is a method that encrypts the password
func (p *Password) Hash() string {
	return p.hash
}

// Compare is a method that compares the password with the encrypted password
func (p *Password) Compare(password string) error {
	if password == "" {
		return errors.New("o campo de senha é obrigatório")
	}

	if err := p.EncryptProvider.Compare(password, p.hash); err != nil {
		return errors.New("credenciais inválidas")
	}

	return nil
}

func (p *Password) ConfirmPassword(passwordToConfirm string) error {
	if passwordToConfirm != p.password {
		return errors.New("as senhas não conferem")
	}

	return nil
}

// Omit is a method that returns an empty string to omit the password
func (p *Password) Omit() string {
	return ""
}
