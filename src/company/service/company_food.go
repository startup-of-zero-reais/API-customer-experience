package service

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/domain"
)

type (
	CompanyFood interface {
		FoodList(slug string) ([]domain.CompanyFood, error)
	}

	CompanyFoodImpl struct {
		Repository data.CompanyFoodsRepository

		logger *providers.LogProvider
	}
)

func NewCompanyFood(repository data.CompanyFoodsRepository, logger *providers.LogProvider) CompanyFood {
	return &CompanyFoodImpl{
		Repository: repository,
		logger:     logger,
	}
}

func (c *CompanyFoodImpl) FoodList(slug string) ([]domain.CompanyFood, error) {
	
	c.logger.WithFields(map[string]interface{}{
		"event": "company_query"
		"company": slug,
	}).Infoln("listing company menu food")
	
	return c.Repository.FoodList()
}
