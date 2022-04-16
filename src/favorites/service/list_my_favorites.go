package service

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/domain"
)

type (
	ListMyFavorites interface {
		List(loggedUsrID string) ([]domain.Favorite, error)
	}

	ListMyFavoritesImpl struct {
		Repository     domain.FavoriteRepository
		MealRepository domain.MealRepository
	}
)

func NewListMyFavorites(repository domain.FavoriteRepository, mealRepository domain.MealRepository) *ListMyFavoritesImpl {
	return &ListMyFavoritesImpl{
		Repository:     repository,
		MealRepository: mealRepository,
	}
}

func (a *ListMyFavoritesImpl) List(loggedUsrID string) ([]domain.Favorite, error) {
	if loggedUsrID == "" {
		return nil, validation.BadRequestError("usuário não encontrado")
	}

	favorites, err := a.Repository.UsrFavorites(loggedUsrID)
	if err != nil {
		return nil, err
	}

	for i, favorite := range favorites {
		meal, err := a.MealRepository.GetMeal(favorite.Meal.ID)
		if err != nil {
			return nil, err
		}

		favorites[i].Meal = *meal
		favorites[i].Company = meal.Company
	}

	return favorites, nil
}
