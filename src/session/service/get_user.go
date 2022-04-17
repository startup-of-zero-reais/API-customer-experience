package service

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"
)

type (
	GetUser interface {
		Find(email string) (*domain.User, error)
	}

	GetUserImpl struct {
		userRepository domain.UserRepository
		logger         *providers.LogProvider
	}
)

func NewGetUser(userRepository domain.UserRepository, logger *providers.LogProvider) GetUser {
	return &GetUserImpl{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (g *GetUserImpl) Find(email string) (*domain.User, error) {
	g.logger.Infoln("getting user", email)
	return g.userRepository.Find(email)
}
