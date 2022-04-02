package handler

import (
	"errors"
	"net/http"

	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/service"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
)

type (
	// Handler struct represents a route handler
	Handler struct {
		response *domain.Response

		createUser *service.CreateUser
	}
)

func NewHandler() Handler {
	userRepository := data.NewUserRepository()

	return Handler{
		response: domain.NewResponse(),

		createUser: service.NewCreateUser(userRepository),
	}
}

func (h *Handler) Post(body string) domain.Response {
	h.response.SetStatusCode(http.StatusCreated)

	err := h.createUser.Execute(body)
	var fv *validation.FieldValidator
	if err != nil && errors.As(err, &fv) {
		h.response.SetStatusCode(http.StatusBadRequest)
		h.response.SetMetadata(err)
	}

	return *h.response
}

func (h *Handler) Get() domain.Response {
	h.response.SetData("Hello World Local!")

	return *h.response
}

func (h *Handler) Put() domain.Response {
	response := domain.NewResponse()

	return *response
}

func (h *Handler) Delete() domain.Response {
	response := domain.NewResponse()

	return *response
}
