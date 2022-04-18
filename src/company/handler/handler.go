package handler

import (
	"net/http"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
)

type (
	Handler struct {
		response *domain.Response

		app *Application

		logger *providers.LogProvider
	}
)

func NewHandler(logger *providers.LogProvider) *Handler {
	return &Handler{
		response: domain.NewResponse(),

		app: NewApplication(logger),

		logger: logger,
	}
}

func (h *Handler) Get(r domain.Request) domain.Response {
	foodList, err := h.app.Queries.CompanyFood.FoodList(r.PathParams["slug"])
	if err != nil {
		return h.response.HandleError(err)
	}

	h.response.SetStatusCode(http.StatusOK)
	h.response.SetData(foodList)

	return *h.response
}
