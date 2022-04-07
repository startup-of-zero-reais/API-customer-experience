package service

import "github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"

type (
	GetUser interface {
		Find(email string) (*domain.User, error)
	}

	GetUserImpl struct {
		userRepository domain.UserRepository
	}
)

func NewGetUser(userRepository domain.UserRepository) GetUser {
	return &GetUserImpl{
		userRepository: userRepository,
	}
}

func (g *GetUserImpl) Find(email string) (*domain.User, error) {
	return g.userRepository.Find(email)
}
