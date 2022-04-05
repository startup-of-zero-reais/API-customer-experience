package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/providers"

	"github.com/google/uuid"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/fields"
)

type (
	// User struct represents a user
	CreateUser struct {
		Repository data.UserRepository
	}

	// User struct represents a user
	User struct {
		ID              string `json:"id,omitempty"`
		Name            string `json:"name,omitempty"`
		Lastname        string `json:"lastname,omitempty"`
		Email           string `json:"email,omitempty"`
		Phone           string `json:"phone,omitempty"`
		Password        string `json:"password,omitempty"`
		ConfirmPassword string `json:"confirm_password,omitempty"`
		Avatar          string `json:"avatar,omitempty"`
	}
)

func NewCreateUser(repository data.UserRepository) *CreateUser {
	return &CreateUser{
		Repository: repository,
	}
}

func (c *CreateUser) Execute(body string) error {
	var receivedUser User
	err := json.Unmarshal([]byte(body), &receivedUser)
	if err != nil {
		return err
	}

	userAlreadyExists, err := c.Repository.FindByEmail(receivedUser.Email)
	if err != nil {
		var notFound *validation.NotFound
		if !errors.As(err, &notFound) {
			return err
		}
	}

	if userAlreadyExists != nil {
		return validation.EntityAlreadyExistsError("usuário com este e-mail já cadastrado")
	}

	user, err := domain.NewUser(
		uuid.NewString(),
		receivedUser.Name,
		receivedUser.Lastname,
		receivedUser.Email,
		receivedUser.Phone,
		receivedUser.Avatar,
		fields.NewPassword(
			providers.NewEncryptProvider(),
			receivedUser.Password,
		),
	)

	if err != nil {
		return err
	}

	err = user.ConfirmPassword(receivedUser.ConfirmPassword)
	if err != nil {
		return err
	}

	err = c.Repository.Save(user)
	if err != nil {
		return err
	}

	log.Println("[USER UID]:", user.ID)

	return nil
}
