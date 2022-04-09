package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	s "github.com/startup-of-zero-reais/API-customer-experience/src/common/service"

	d "github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"
)

type (
	Handler struct {
		response *domain.Response

		app *Application

		jwtService s.JwtService
	}
)

func NewHandler() Handler {
	jwt := s.NewJwtService()

	return Handler{
		response: domain.NewResponse(),

		app: NewApplication(
			jwt,
		),

		jwtService: jwt,
	}
}

func (h *Handler) SignIn(r domain.Request) domain.Response {
	var authInput d.AuthInput
	err := json.Unmarshal([]byte(r.Body), &authInput)
	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	session, err := h.app.Commands.SignIn.SignIn(authInput.Email, authInput.Password)
	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	tokenCookie := fmt.Sprintf("usess=%s", session.SessionToken)
	h.response.MultiValueHeaders["Set-Cookie"] = []string{tokenCookie}
	h.response.MultiValueHeaders["set-Cookie"] = []string{tokenCookie}
	h.response.Cookies = []string{tokenCookie}

	log.Printf("\n\n\nAuth Token:\n\n%s\n\n", session.SessionToken)

	return *h.response
}

func (h *Handler) SignOut(r domain.Request) domain.Response {
	isValid, err := h.jwtService.ValidateToken(r.Cookies["usess"])

	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	if !isValid {
		h.response.SetStatusCode(http.StatusBadRequest)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	id := h.jwtService.DecodedToken("id").(string)
	err = h.app.SignOut.ClearSession(id)
	if err != nil {
		h.response.SetStatusCode(http.StatusBadRequest)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	h.response.SetStatusCode(http.StatusNoContent)

	return *h.response
}

func (h *Handler) RecoverPassword(r domain.Request) domain.Response {
	return *h.response
}
