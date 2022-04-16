package service

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/domain"
)

type (
	RemoveFromFavorite interface {
		Favorite(loggedUsrId, favoriteID string) error
	}

	RemoveFromFavoriteImpl struct {
		Repository domain.FavoriteRepository
	}
)

func NewRemoveFromFavorite(repository domain.FavoriteRepository) *RemoveFromFavoriteImpl {
	return &RemoveFromFavoriteImpl{
		Repository: repository,
	}
}

func (a *RemoveFromFavoriteImpl) Favorite(loggedUsrId, favoriteID string) error {
	if favoriteID == "" {
		return validation.BadRequestError("o prato favorito deve ser informado")
	}

	return a.Repository.Delete(loggedUsrId, favoriteID)
}
