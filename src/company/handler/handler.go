package handler

import (
	"net/http"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
)

type (
	Handler struct {
		response *domain.Response

		app *Application
	}
)

func NewHandler() *Handler {
	return &Handler{
		response: domain.NewResponse(),

		app: NewApplication(),
	}
}

func (h *Handler) Get(r domain.Request) domain.Response {
	foodList, err := h.app.Queries.CompanyFood.FoodList()
	if err != nil {
		h.response.SetStatusCode(http.StatusBadGateway)
		h.response.SetMetadata(
			map[string]string{"error": err.Error()},
		)

		return *h.response
	}

	h.response.SetStatusCode(http.StatusOK)
	h.response.SetData(foodList)

	return *h.response
}
