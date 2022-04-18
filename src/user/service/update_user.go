package service

import (
	"encoding/json"

	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/fields"
)

type (
	// UpdateUser is a interface to update a user
	UpdateUser interface {
		Update(id, email, body string) error
	}

	// UpdateUserImpl is a implementation of UpdateUser
	UpdateUserImpl struct {
		Repository data.UserRepository

		logger *providers.LogProvider
	}

	// User struct represents a user
	UpdateUserInput struct {
		Name               string `json:"name,omitempty"`
		Lastname           string `json:"lastname,omitempty"`
		Email              string `json:"email,omitempty"`
		Phone              string `json:"phone,omitempty"`
		Password           string `json:"password,omitempty"`
		ConfirmPassword    string `json:"confirm_password,omitempty"`
		NewPassword        string `json:"new_password,omitempty"`
		ConfirmNewPassword string `json:"confirm_new_password,omitempty"`
		Avatar             string `json:"avatar,omitempty"`
	}
)

func NewUpdateUser(repository data.UserRepository, logger *providers.LogProvider) UpdateUser {
	return &UpdateUserImpl{
		Repository: repository,

		logger: logger,
	}
}

func (u *UpdateUserImpl) Update(id, email, body string) error {
	existsUser, err := u.Repository.FindByUser(id, email)
	if err != nil {
		return err
	}

	if existsUser == nil {
		return validation.NotFoundError("usuário não encontrado")
	}

	var input UpdateUserInput
	err = json.Unmarshal([]byte(body), &input)
	if err != nil {
		return err
	}

	return u.Repository.Update(id, email, func(user *domain.User) (*domain.User, error) {
		u.logger.Debugf("[INPUT]: %+v\n", input)
		// confirm if is users password
		err := user.Password.PassToHash().Compare(input.Password)
		if err != nil {
			return nil, err
		}

		// re instantiate user password with old values
		user.Password = fields.NewPassword(
			providers.NewEncryptProvider(),
			input.Password,
		)

		// confirm if passwords match
		err = user.ConfirmPassword(input.ConfirmPassword)
		if err != nil {
			return nil, err
		}

		err = u.reflectUser(user, input)
		if err != nil {
			return nil, err
		}

		logFields := map[string]interface{}{
			"user_id": user.ID,
			"event":   d.UserUpdated,
		}

		if input.NewPassword != "" {
			u.logger.WithFields(logFields).Infoln(user.ToString(), "user updated with password")
		} else {
			u.logger.WithFields(logFields).Infoln(user.ToString(), "user updated")
		}

		return user, nil
	})
}

func (u *UpdateUserImpl) reflectUser(user *domain.User, input UpdateUserInput) error {
	if input.Name != "" {
		user.Name = input.Name
	}

	if input.Lastname != "" {
		user.Lastname = input.Lastname
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	if input.Phone != "" {
		user.Phone = input.Phone
	}

	// when NewPassword is not empty, it means that the user wants to change the password
	if input.NewPassword != "" {
		user.Password = fields.NewPassword(
			providers.NewEncryptProvider(),
			input.NewPassword,
		)

		err := user.ConfirmPassword(input.ConfirmNewPassword)
		if err != nil {
			return err
		}
	}

	if input.Avatar != "" {
		user.Avatar = input.Avatar
	}

	return nil
}
