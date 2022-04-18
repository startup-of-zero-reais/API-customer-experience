package service

import (
	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/domain"
)

type (
	AddToFavorite interface {
		Meal(loggedUsrID, mealID string) error
	}

	AddToFavoriteImpl struct {
		Repository     domain.FavoriteRepository
		MealRepository domain.MealRepository

		logger *providers.LogProvider
	}
)

func NewAddToFavorite(repository domain.FavoriteRepository, mealRepository domain.MealRepository, logger *providers.LogProvider) *AddToFavoriteImpl {
	return &AddToFavoriteImpl{
		Repository:     repository,
		MealRepository: mealRepository,
		logger:         logger,
	}
}

func (a *AddToFavoriteImpl) Meal(loggedUsrID, mealID string) error {
	if loggedUsrID == "" {
		return validation.BadRequestError("usuário não logado")
	}

	if mealID == "" {
		return validation.BadRequestError("o prato favorito deve ser informado")
	}

	usrFavs, err := a.Repository.UsrFavorites(loggedUsrID)
	if err != nil {
		return err
	}

	for _, usrFav := range usrFavs {
		if usrFav.Meal.ID == mealID {
			return validation.BadRequestError("prato já adicionado aos favoritos")
		}
	}

	meal, err := a.MealRepository.GetMeal(mealID)
	if err != nil {
		return validation.NotFoundError("prato não encontrado")
	}

	favorite, err := domain.NewFavorite("", loggedUsrID, meal.Company, meal)
	if err != nil {
		return err
	}

	a.logger.WithFields(map[string]interface{}{
		"user_id": loggedUsrID,
		"event":   d.FavoriteAdded,
	}).Infoln("adding favorite to user")

	return a.Repository.Add(favorite)
}
