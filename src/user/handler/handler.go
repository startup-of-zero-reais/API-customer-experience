package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/service"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	s "github.com/startup-of-zero-reais/API-customer-experience/src/common/service"
)

type (
	// Handler struct represents a route handler
	Handler struct {
		response *domain.Response

		app *Application
		*providers.LogProvider

		jwtService s.JwtService
	}
)

func NewHandler() Handler {
	logger := providers.NewLogProvider()

	return Handler{
		response: domain.NewResponse(),

		app:         NewApplication(logger),
		LogProvider: logger,

		jwtService: s.NewJwtService(),
	}
}

func (h *Handler) validateAuth(r domain.Request) error {
	authorization := r.Headers["Authorization"]
	if authorization == "" {
		authorization = r.Cookies["usess"]
	}

	if authorization == "" {
		return errors.New("sessão expirada ou inválida")
	}

	_, err := h.jwtService.ValidateToken(authorization)
	return err
}

func (h *Handler) Post(r domain.Request) domain.Response {
	h.response.SetStatusCode(http.StatusCreated)

	err := h.app.CreateUser.Execute(r.Body)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	return *h.response
}

func (h *Handler) Get(r domain.Request) domain.Response {
	err := h.validateAuth(r)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	id := h.jwtService.DecodedToken("id").(string)
	email := h.jwtService.DecodedToken("email").(string)

	user, err := h.app.GetUser.Execute(id, email)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	h.response.SetStatusCode(http.StatusOK)
	h.response.SetData(user)
	h.LogResponse(*h.response)

	return *h.response
}

func (h *Handler) Put(r domain.Request) domain.Response {
	err := h.validateAuth(r)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	id := h.jwtService.DecodedToken("id").(string)
	email := h.jwtService.DecodedToken("email").(string)

	err = h.app.UpdateUser.Update(id, email, r.Body)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	return *h.response
}

func (h *Handler) Delete(r domain.Request) domain.Response {
	err := h.validateAuth(r)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	id := h.jwtService.DecodedToken("id").(string)
	email := h.jwtService.DecodedToken("email").(string)

	var user service.User
	err = json.Unmarshal([]byte(r.Body), &user)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	err = h.app.DeleteUser.Execute(id, email, user.Password, user.ConfirmPassword)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	h.response.SetStatusCode(http.StatusNoContent)
	return *h.response
}
