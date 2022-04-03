package service

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
)

type (
	GetUser struct {
		Repository data.UserRepository
	}
)

func NewGetUser(repository data.UserRepository) *GetUser {
	return &GetUser{
		Repository: repository,
	}
}

func (c *GetUser) Execute(id, email string) (*data.UserModel, error) {
	return c.Repository.FindByUser(id, email)
}
