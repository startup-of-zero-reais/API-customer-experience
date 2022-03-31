package domain

import (
	"errors"
)

type (
	Address struct {
		ID           string
		Street       string
		City         string
		State        string
		ZipCode      string
		Neighborhood string
	}

	// User struct represents a user
	User struct {
		ID        string
		Name      string
		Lastname  string
		Email     string
		Phone     string
		Password  Password
		Avatar    string
		Addresses []Address
	}
)

func NewUser(id, name, lastname, email, phone, avatar string, password Password) *User {
	if err := requiredStringField(id); err != nil {
		panic(err)
	}

	if err := requiredStringField(name); err != nil {
		panic(err)
	}

	if err := requiredStringField(lastname); err != nil {
		panic(err)
	}

	if err := requiredStringField(email); err != nil {
		panic(err)
	}

	if err := requiredStringField(phone); err != nil {
		panic(err)
	}

	if err := requiredStringField(avatar); err != nil {
		panic(err)
	}

	if !password.IsValid() {
		panic(errors.New("o campo de senha é obrigatório"))
	}

	return &User{
		ID:       id,
		Name:     name,
		Lastname: lastname,
		Email:    email,
		Phone:    phone,
		Avatar:   avatar,
		Password: password.Hash(),
	}
}
