package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	s "github.com/startup-of-zero-reais/API-customer-experience/src/common/service"

	d "github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"
)

type (
	Handler struct {
		response *domain.Response

		app *Application
		*providers.LogProvider

		jwtService s.JwtService
	}
)

func NewHandler(l *providers.LogProvider) Handler {
	jwt := s.NewJwtService()

	return Handler{
		response: domain.NewResponse(),

		app: NewApplication(
			jwt,
			l,
		),
		LogProvider: l,

		jwtService: jwt,
	}
}

func (h *Handler) SignIn(r domain.Request) domain.Response {
	var authInput d.AuthInput
	err := json.Unmarshal([]byte(r.Body), &authInput)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	session, err := h.app.Commands.SignIn.SignIn(authInput.Email, authInput.Password)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	tokenCookie := fmt.Sprintf("usess=%s; expires=%s;", session.SessionToken, time.Unix(session.ExpiresIn, 0).Format(time.RFC1123))
	if os.Getenv("ENVIRONMENT") == "production" {
		tokenCookie = fmt.Sprintf("usess=%s; domain=zero-reais-lab.cloud; expires=%s;", session.SessionToken, time.Unix(session.ExpiresIn, 0).Format(time.RFC1123))
	}

	h.response.Headers["X-Auth-Token"] = session.SessionToken
	h.response.Cookies = []string{tokenCookie}

	h.LogResponse(*h.response)

	return *h.response
}

func (h *Handler) SignOut(r domain.Request) domain.Response {
	isValid, err := h.jwtService.ValidateToken(r.Cookies["usess"])

	if err != nil {
		if !strings.Contains(err.Error(), "token is expired by") {
			res := h.response.HandleError(err)
			h.LogResponse(res)
			return res
		}
	}

	if !isValid {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	id := h.jwtService.DecodedToken("id").(string)
	err = h.app.SignOut.ClearSession(id)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	h.response.SetStatusCode(http.StatusNoContent)

	h.LogResponse(*h.response)

	return *h.response
}

func (h *Handler) RecoverPassword(r domain.Request) domain.Response {
	var authInput d.AuthInput
	err := json.Unmarshal([]byte(r.Body), &authInput)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	err = h.app.RecoverPassword.SendOTP(authInput.Email)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	h.LogResponse(*h.response)

	return *h.response
}

func (h *Handler) ResetPassword(r domain.Request) domain.Response {
	h.response.SetStatusCode(http.StatusNoContent)

	var resetPassInput d.ResetPassInput
	err := json.Unmarshal([]byte(r.Body), &resetPassInput)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	if resetPassInput.Password != resetPassInput.ConfirmPassword {
		h.response.SetStatusCode(http.StatusBadRequest)
		h.response.SetMetadata(map[string]interface{}{"error": "senhas não conferem"})

		return *h.response
	}

	err = h.app.RecoverPassword.ResetPassword(resetPassInput.OTP, resetPassInput.Password)
	if err != nil {
		res := h.response.HandleError(err)
		h.LogResponse(res)
		return res
	}

	return *h.response
}
