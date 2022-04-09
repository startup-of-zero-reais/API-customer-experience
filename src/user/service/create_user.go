package service

import (
	"encoding/json"
	"errors"
	"log"

	dt "github.com/startup-of-zero-reais/API-customer-experience/src/common/data"
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

		EventEmitter d.EventEmitter
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

func NewCreateUser(repository data.UserRepository) *CreateUserImpl {
	return &CreateUserImpl{
		Repository: repository,

		EventEmitter: d.NewEventEmitter(
			dt.NewEventRepository(),
		),
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

	log.Println("[USER UID]:", user.ID)
	return c.EventEmitter.Emit(user.ID, d.UserCreated, user)
}
