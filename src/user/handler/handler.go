package handler

import (
	"net/http"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
)

type (
	// Handler struct represents a route handler
	Handler struct{}
)

func (h *Handler) Post(body string) domain.Response {
	response := domain.NewResponse()
	response.SetStatusCode(http.StatusCreated)
	response.SetData(body)

	return *response
}

func (h *Handler) Get() domain.Response {
	response := domain.NewResponse()
	response.SetData("Hello World Local!")

	return *response
}

func (h *Handler) Put() domain.Response {
	response := domain.NewResponse()

	return *response
}

func (h *Handler) Delete() domain.Response {
	response := domain.NewResponse()

	return *response
}
