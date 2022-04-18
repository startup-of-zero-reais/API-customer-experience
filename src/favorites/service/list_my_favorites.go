package service

import (
	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
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
		logger         *providers.LogProvider
	}
)

func NewListMyFavorites(repository domain.FavoriteRepository, mealRepository domain.MealRepository, logger *providers.LogProvider) *ListMyFavoritesImpl {
	return &ListMyFavoritesImpl{
		Repository:     repository,
		MealRepository: mealRepository,
		logger:         logger,
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

	a.logger.WithFields(map[string]interface{}{
		"user_id": loggedUsrID,
		"event":   d.FavoriteQuery,
	}).Infoln("listing favorites from user")

	return favorites, nil
}
