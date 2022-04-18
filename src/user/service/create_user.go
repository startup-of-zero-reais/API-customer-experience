package service

import (
	"encoding/json"
	"errors"

	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"

	"github.com/google/uuid"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/fields"
)

type (
	CreateUser interface {
		Execute(body string) error
	}

	// User struct represents a user
	CreateUserImpl struct {
		Repository data.UserRepository

		logger *providers.LogProvider
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

func NewCreateUser(repository data.UserRepository, logger *providers.LogProvider) *CreateUserImpl {
	return &CreateUserImpl{
		Repository: repository,

		logger: logger,
	}
}

func (c *CreateUserImpl) Execute(body string) error {
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

	c.logger.WithFields(map[string]interface{}{
		"user_id": user.ID,
		"event":   d.UserCreated,
	}).Infoln(user.ToString())

	return nil
}
