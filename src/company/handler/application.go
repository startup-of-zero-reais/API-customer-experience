package handler

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/service"
)

type (
	Queries struct {
		CompanyFood service.CompanyFood
	}

	Application struct {
		Queries Queries

		logger *providers.LogProvider
	}
)

func NewApplication(logger *providers.LogProvider) *Application {
	cRepository := data.NewCompanyFoodsRepository()

	return &Application{
		Queries: Queries{
			CompanyFood: service.NewCompanyFood(
				cRepository,
				logger,
			),
		},
	}
}
