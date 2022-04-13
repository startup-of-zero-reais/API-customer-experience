package service

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/domain"
)

type (
	CompanyFood interface {
		FoodList() ([]domain.CompanyFood, error)
	}

	CompanyFoodImpl struct {
		Repository data.CompanyFoodsRepository
	}
)

func NewCompanyFood(repository data.CompanyFoodsRepository) CompanyFood {
	return &CompanyFoodImpl{
		Repository: repository,
	}
}

func (c *CompanyFoodImpl) FoodList() ([]domain.CompanyFood, error) {
	return c.Repository.FoodList()
}
