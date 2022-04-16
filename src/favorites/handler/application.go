package handler

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/service"
)

type (
	Commands struct {
		AddToFavorite      service.AddToFavorite
		RemoveFromFavorite service.RemoveFromFavorite
	}

	Queries struct {
		ListMyFavorites service.ListMyFavorites
	}

	Application struct {
		Commands Commands
		Queries  Queries
	}
)

func NewApplication() *Application {
	favoritesRepository := data.NewFavoritesRepository()
	mealRepository := data.NewMealRepository()

	return &Application{
		Commands: Commands{
			AddToFavorite: service.NewAddToFavorite(
				favoritesRepository,
				mealRepository,
			),
			RemoveFromFavorite: service.NewRemoveFromFavorite(favoritesRepository),
		},
		Queries: Queries{
			ListMyFavorites: service.NewListMyFavorites(
				favoritesRepository,
				mealRepository,
			),
		},
	}
}
