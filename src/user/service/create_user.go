package service

import (
	"encoding/json"

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

	user, err := domain.NewUser(
		uuid.NewString(),
		receivedUser.Name,
		receivedUser.Lastname,
		receivedUser.Email,
		receivedUser.Phone,
		receivedUser.Avatar,
		fields.NewPassword(
			NewEncryptProvider(),
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

	user.Password.Encrypt()

	err = c.Repository.Save(user)
	if err != nil {
		return err
	}

	return nil
}
