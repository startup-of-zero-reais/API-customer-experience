package service

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
)

type (
	GetUser interface {
		Execute(id, email string) (*data.UserModel, error)
	}

	GetUserImpl struct {
		Repository data.UserRepository
		logger     *providers.LogProvider
	}
)

func NewGetUser(repository data.UserRepository, logger *providers.LogProvider) *GetUserImpl {
	return &GetUserImpl{
		Repository: repository,
		logger:     logger,
	}
}

func (c *GetUserImpl) Execute(id, email string) (*data.UserModel, error) {
	c.logger.WithFields(map[string]interface{}{
		"user_id": id,
		"event":   "query_user",
	}).Infoln("query command user email", email)
	return c.Repository.FindByUser(id, email)
}
