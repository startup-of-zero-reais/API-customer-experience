package service

import (
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
	}
)

func NewAddToFavorite(repository domain.FavoriteRepository, mealRepository domain.MealRepository) *AddToFavoriteImpl {
	return &AddToFavoriteImpl{
		Repository:     repository,
		MealRepository: mealRepository,
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

	return a.Repository.Add(favorite)
}
