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
		getUser    *service.GetUser
	}
)

func NewHandler() Handler {
	userRepository := data.NewUserRepository()

	return Handler{
		response: domain.NewResponse(),

		createUser: service.NewCreateUser(userRepository),
		getUser:    service.NewGetUser(userRepository),
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

func (h *Handler) Get(headers map[string]string) domain.Response {
	if len(headers) <= 0 {
		h.response.SetStatusCode(http.StatusBadRequest)
		h.response.SetMetadata(map[string]string{"error": "Header 'Authorization' é obrigatório"})

		return *h.response
	}

	id := headers["User-Id"]
	email := headers["User-Email"]
	if id == "" || email == "" {
		h.response.SetStatusCode(http.StatusBadRequest)
		h.response.SetMetadata(map[string]string{"error": "Header 'Authorization' é obrigatório"})

		return *h.response
	}

	user, err := h.getUser.Execute(id, email)
	if err != nil {
		h.response.SetStatusCode(http.StatusNotFound)
		h.response.SetMetadata(map[string]string{"error": err.Error()})
	} else {
		h.response.SetStatusCode(http.StatusOK)
		h.response.SetData(user)
	}

	return *h.response
}

func (h *Handler) Put() domain.Response {
	h.response.SetStatusCode(http.StatusCreated)

	return *h.response
}

func (h *Handler) Delete() domain.Response {
	h.response.SetStatusCode(http.StatusCreated)

	return *h.response
}
