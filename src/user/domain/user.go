package domain

import (
	"encoding/json"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/fields"
)

type (
	// User struct represents a user
	User struct {
		ID        string            `json:"id,omitempty"`
		Name      string            `json:"name"`
		Lastname  string            `json:"lastname"`
		Email     string            `json:"email"`
		Phone     string            `json:"phone"`
		Password  *fields.Password  `json:"password,omitempty"`
		Avatar    string            `json:"avatar,omitempty"`
		Addresses *[]fields.Address `json:"addresses,omitempty"`
	}
)

// NewUser is the constructor of User
func NewUser(id, name, lastname, email, phone, avatar string, password *fields.Password) (*User, error) {
	fieldValidator := validation.NewFieldValidator()

	if err := validation.RequiredStringField(id, "id"); err != nil {
		fieldValidator.AddError("id", err.Error())
	}

	if err := validation.RequiredStringField(name, "nome"); err != nil {
		fieldValidator.AddError("name", err.Error())
	}

	if err := validation.RequiredStringField(lastname, "sobrenome"); err != nil {
		fieldValidator.AddError("lastname", err.Error())
	}

	if err := validation.RequiredStringField(email, "e-mail"); err != nil {
		fieldValidator.AddError("email", err.Error())
	}

	if err := validation.RequiredStringField(phone, "telefone"); err != nil {
		fieldValidator.AddError("phone", err.Error())
	}

	if errs := password.IsValid(); len(errs) > 0 {
		for _, errString := range errs {
			fieldValidator.AddError("password", errString)
		}
	}

	if fieldValidator.HasErrors() {
		return nil, fieldValidator
	}

	return &User{
		ID:        id,
		Name:      name,
		Lastname:  lastname,
		Email:     email,
		Phone:     phone,
		Avatar:    avatar,
		Password:  password,
		Addresses: &[]fields.Address{},
	}, nil
}

func (u *User) ConfirmPassword(password string) error {
	fieldValidator := validation.NewFieldValidator()

	if password == "" {
		fieldValidator.AddError("confirm_password", "o campo de confirmação de senha é obrigatório")
	}

	if err := u.Password.ConfirmPassword(password); err != nil {
		fieldValidator.AddError("confirm_password", err.Error())
	}

	if fieldValidator.HasErrors() {
		return fieldValidator
	}

	return nil
}

func (u *User) ToString() string {
	bytes, _ := json.Marshal(u)

	return string(bytes)
}
