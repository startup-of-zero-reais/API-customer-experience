package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	s "github.com/startup-of-zero-reais/API-customer-experience/src/common/service"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"

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
		var notFound *validation.NotFound
		if errors.As(err, &notFound) {
			h.response.SetStatusCode(http.StatusNotFound)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

			return *h.response
		}

		var unauthorized *validation.Unauthorized
		if errors.As(err, &unauthorized) {
			h.response.SetStatusCode(http.StatusUnauthorized)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

			return *h.response
		}

		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	tokenCookie := fmt.Sprintf("usess=%s; domain=zero-reais-lab.cloud; expires=%s;", session.SessionToken, time.Unix(session.ExpiresIn, 0).Format(time.RFC1123))
	h.response.Headers["set-cookie"] = tokenCookie
	h.response.Headers["x-auth-token"] = session.SessionToken
	h.response.Cookies = []string{tokenCookie}

	log.Printf("\n\n\nSet-Cookie:\n\n%s\n\n", tokenCookie)
	log.Printf("\n\n\nAuth Token:\n\n%s\n\n", session.SessionToken)
	log.Printf("\n\n\nHeaders:\n\n%+v\n\n", h.response.Headers)

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
	var authInput d.AuthInput
	err := json.Unmarshal([]byte(r.Body), &authInput)
	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	err = h.app.RecoverPassword.SendOTP(authInput.Email)
	if err != nil {
		var notFound *validation.NotFound
		if errors.As(err, &notFound) {
			h.response.SetStatusCode(http.StatusNotFound)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

			return *h.response
		}

		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	return *h.response
}

func (h *Handler) ResetPassword(r domain.Request) domain.Response {
	h.response.SetStatusCode(http.StatusNoContent)

	var resetPassInput d.ResetPassInput
	err := json.Unmarshal([]byte(r.Body), &resetPassInput)
	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	if resetPassInput.Password != resetPassInput.ConfirmPassword {
		h.response.SetStatusCode(http.StatusBadRequest)
		h.response.SetMetadata(map[string]interface{}{"error": "senhas não conferem"})

		return *h.response
	}

	err = h.app.RecoverPassword.ResetPassword(resetPassInput.OTP, resetPassInput.Password)
	if err != nil {
		var unauthorized *validation.Unauthorized
		if errors.As(err, &unauthorized) {
			h.response.SetStatusCode(http.StatusUnauthorized)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

			return *h.response
		}

		var notFound *validation.NotFound
		if errors.As(err, &notFound) {
			h.response.SetStatusCode(http.StatusNotFound)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

			return *h.response
		}

		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	return *h.response
}
