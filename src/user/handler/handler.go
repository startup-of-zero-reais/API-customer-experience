package handler

import (
	"errors"
	"net/http"

	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain/validation"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	s "github.com/startup-of-zero-reais/API-customer-experience/src/common/service"
)

type (
	// Handler struct represents a route handler
	Handler struct {
		response *domain.Response

		app *Application

		jwtService s.JwtService
	}

	Request struct {
		Headers map[string]string
		Cookies map[string]string
		Body    string
	}
)

func NewHandler() Handler {
	return Handler{
		response: domain.NewResponse(),

		app: NewApplication(),

		jwtService: s.NewJwtService(),
	}
}

func (h *Handler) validateAuth(r Request) error {
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

func (h *Handler) Post(r Request) domain.Response {
	h.response.SetStatusCode(http.StatusCreated)

	err := h.app.CreateUser.Execute(r.Body)
	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		var fv *validation.FieldValidator
		if errors.As(err, &fv) {
			h.response.SetStatusCode(http.StatusBadRequest)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})
		}

		var alreadyExists *validation.EntityAlreadyExists
		if errors.As(err, &alreadyExists) {
			h.response.SetStatusCode(http.StatusConflict)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})
		}
	}

	return *h.response
}

func (h *Handler) Get(r Request) domain.Response {
	// log.Println(
	// 	h.jwtService.GenerateToken(
	// 		"cd305516-9f35-461c-8af4-6f8a53278398",
	// 		"john@doe.com",
	// 	),
	// )
	err := h.validateAuth(r)

	if err != nil {
		h.response.SetStatusCode(http.StatusUnauthorized)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	id := h.jwtService.DecodedToken("id").(string)
	email := h.jwtService.DecodedToken("email").(string)

	user, err := h.app.GetUser.Execute(id, email)
	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		var notFound *validation.NotFound
		if errors.As(err, &notFound) {
			h.response.SetStatusCode(http.StatusNotFound)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})
		}
	} else {
		h.response.SetStatusCode(http.StatusOK)
		h.response.SetData(user)
	}

	return *h.response
}

func (h *Handler) Put(r Request) domain.Response {
	err := h.validateAuth(r)
	if err != nil {
		h.response.SetStatusCode(http.StatusUnauthorized)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		return *h.response
	}

	id := h.jwtService.DecodedToken("id").(string)
	email := h.jwtService.DecodedToken("email").(string)

	err = h.app.UpdateUser.Update(id, email, r.Body)
	if err != nil {
		h.response.SetStatusCode(http.StatusInternalServerError)
		h.response.SetMetadata(map[string]interface{}{"error": err.Error()})

		var notFound *validation.NotFound
		if errors.As(err, &notFound) {
			h.response.SetStatusCode(http.StatusNotFound)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})
		}

		var fv *validation.FieldValidator
		if errors.As(err, &fv) {
			h.response.SetStatusCode(http.StatusBadRequest)
			h.response.SetMetadata(err)
		}

		var alreadyExists *validation.EntityAlreadyExists
		if errors.As(err, &alreadyExists) {
			h.response.SetStatusCode(http.StatusConflict)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})
		}

		var unauthorized *validation.Unauthorized
		if errors.As(err, &unauthorized) {
			h.response.SetStatusCode(http.StatusUnauthorized)
			h.response.SetMetadata(map[string]interface{}{"error": err.Error()})
		}
	}

	return *h.response
}

func (h *Handler) Delete() domain.Response {
	h.response.SetStatusCode(http.StatusCreated)

	return *h.response
}
