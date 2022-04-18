package service

import (
	"errors"

	cd "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/fields"
)

type (
	// DeleteUser interface represents a delete user service
	DeleteUser interface {
		Execute(id, email, password, confirmPassword string) error
	}

	// DeleteUserImpl struct represents a delete user service implementation
	DeleteUserImpl struct {
		Repository data.UserRepository

		logger *providers.LogProvider
	}
)

// NewDeleteUser returns a new instance of DeleteUser
func NewDeleteUser(repository data.UserRepository, logger *providers.LogProvider) DeleteUser {
	return &DeleteUserImpl{
		Repository: repository,

		logger: logger,
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

	d.logger.WithFields(map[string]interface{}{
		"user_id": id,
		"event":   cd.UserDeleted,
	}).Infoln(email)

	return d.Repository.Delete(id, email)
}
