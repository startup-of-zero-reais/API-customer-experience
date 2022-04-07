package service

import (
	"errors"

	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/fields"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/providers"
)

type (
	// DeleteUser interface represents a delete user service
	DeleteUser interface {
		Execute(id, email, password, confirmPassword string) error
	}

	// DeleteUserImpl struct represents a delete user service implementation
	DeleteUserImpl struct {
		Repository data.UserRepository

		EventEmitter domain.EventEmitter
	}
)

// NewDeleteUser returns a new instance of DeleteUser
func NewDeleteUser(repository data.UserRepository) DeleteUser {
	return &DeleteUserImpl{
		Repository: repository,

		EventEmitter: domain.NewEventEmitter(
			data.NewEventRepository(),
		),
	}
}

// Execute deletes a user
func (d *DeleteUserImpl) Execute(id, email, password, confirmPassword string) error {
	userModel, err := d.Repository.FindByUser(id, email)
	if err != nil {
		return err
	}

	user, err := domain.NewUser(
		userModel.ID,
		userModel.Name,
		userModel.Lastname,
		userModel.Email,
		userModel.Phone,
		userModel.Avatar,
		fields.NewPassword(
			providers.NewEncryptProvider(),
			userModel.Password,
		).PassToHash(),
	)
	if err != nil {
		var fv *validation.FieldValidator
		if errors.As(err, &fv) && password != confirmPassword {
			fv.AddError("confirm_password", "as senhas não conferem")
		}

		return fv
	}

	err = user.ConfirmPassword(password)
	if err != nil {
		return err
	}

	err = d.EventEmitter.Emit(id, domain.UserDeleted, email)
	if err != nil {
		return err
	}

	return d.Repository.Delete(id, email)
}
