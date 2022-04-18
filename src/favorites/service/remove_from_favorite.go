package service

import (
	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/domain"
)

type (
	RemoveFromFavorite interface {
		Favorite(loggedUsrId, favoriteID string) error
	}

	RemoveFromFavoriteImpl struct {
		Repository domain.FavoriteRepository
		logger     *providers.LogProvider
	}
)

func NewRemoveFromFavorite(repository domain.FavoriteRepository, logger *providers.LogProvider) *RemoveFromFavoriteImpl {
	return &RemoveFromFavoriteImpl{
		Repository: repository,
		logger:     logger,
	}
}

func (a *RemoveFromFavoriteImpl) Favorite(loggedUsrId, favoriteID string) error {
	if favoriteID == "" {
		return validation.BadRequestError("o prato favorito deve ser informado")
	}

	a.logger.WithFields(map[string]interface{}{
		"user_id": loggedUsrId,
		"event":   d.FavoriteRemoved,
	}).Infoln("removing", favoriteID, "from favorites")

	return a.Repository.Delete(loggedUsrId, favoriteID)
}
