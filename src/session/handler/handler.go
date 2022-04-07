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
	return Handler{
		response: domain.NewResponse(),

		app: NewApplication(),

		jwtService: s.NewJwtService(),
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

	user, err := h.app.Queries.GetUser.Find(authInput.Email)
	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})
	}

	tokenCookie := fmt.Sprintf("usess=%s", token)
	h.response.MultiValueHeaders["Set-Cookie"] = []string{tokenCookie}
	h.response.MultiValueHeaders["set-Cookie"] = []string{tokenCookie}
	h.response.Cookies = []string{tokenCookie}

	log.Printf("\n\n\nAuth Token:\n\n%s\n\n", token)

	return *h.response
}

func (h *Handler) SignOut(r domain.Request) domain.Response {
	return *h.response
}

func (h *Handler) RecoverPassword(r domain.Request) domain.Response {
	return *h.response
}
