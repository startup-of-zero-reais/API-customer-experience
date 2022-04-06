package service

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
)

type (
	GetUser interface {
		Execute(id, email string) (*data.UserModel, error)
	}

	GetUserImpl struct {
		Repository data.UserRepository
	}
)

func NewGetUser(repository data.UserRepository) *GetUserImpl {
	return &GetUserImpl{
		Repository: repository,
	}
}

func (c *GetUserImpl) Execute(id, email string) (*data.UserModel, error) {
	return c.Repository.FindByUser(id, email)
}
