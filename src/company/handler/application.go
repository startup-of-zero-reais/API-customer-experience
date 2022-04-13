package handler

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/service"
)

type (
	Queries struct {
		CompanyFood service.CompanyFood
	}

	Application struct {
		Queries Queries
	}
)

func NewApplication() *Application {
	cRepository := data.NewCompanyFoodsRepository()

	return &Application{
		Queries: Queries{
			CompanyFood: service.NewCompanyFood(
				cRepository,
			),
		},
	}
}
