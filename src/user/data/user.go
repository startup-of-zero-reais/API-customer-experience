package data

import "github.com/startup-of-zero-reais/API-customer-experience/src/user/domain"

type (
	// UserRepository is a interface to access the user data
	UserRepository interface {
		Save(user *domain.User) error
	}

	// UserRepositoryImpl is a implementation of UserRepository
	UserRepositoryImpl struct{}
)

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Save(user *domain.User) error {
	return nil
}
