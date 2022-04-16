package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/service"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
)

type (
	Handler struct {
		response *domain.Response
		app      *Application
		jwt      service.JwtService
	}
)

func NewHandler() *Handler {
	return &Handler{
		response: domain.NewResponse(),
		app:      NewApplication(),
		jwt:      service.NewJwtService(),
	}
}

func (h *Handler) AddToFavorite(r domain.Request) domain.Response {
	err := h.validateAuth(r)
	if err != nil {
		return h.ProperlyError(err)
	}

	loggedUsrId := h.jwt.DecodedToken("id").(string)

	mealInput := struct {
		ID string `json:"meal"`
	}{}

	err = json.Unmarshal([]byte(r.Body), &mealInput)
	if err != nil {
		return h.ProperlyError(err)
	}

	if mealInput.ID == "" {
		return h.ProperlyError(
			validation.BadRequestError("o prato favorito deve ser informado"),
		)
	}

	err = h.app.Commands.AddToFavorite.Meal(loggedUsrId, mealInput.ID)
	if err != nil {
		return h.ProperlyError(err)
	}

	return *h.response
}

func (h *Handler) RemoveFavorite(r domain.Request) domain.Response {
	err := h.validateAuth(r)
	if err != nil {
		return h.ProperlyError(err)
	}

	favorite := struct {
		ID string `json:"favorite"`
	}{}

	err = json.Unmarshal([]byte(r.Body), &favorite)
	if err != nil {
		return h.ProperlyError(err)
	}

	if favorite.ID == "" {
		return h.ProperlyError(
			validation.BadRequestError("o favorito deve ser informado"),
		)
	}

	err = h.app.Commands.RemoveFromFavorite.Favorite(favorite.ID)
	if err != nil {
		return h.ProperlyError(err)
	}

	return *h.response
}

func (h *Handler) MyFavorites(r domain.Request) domain.Response {
	err := h.validateAuth(r)
	if err != nil {
		return h.ProperlyError(err)
	}

	loggedUsrId := h.jwt.DecodedToken("id").(string)

	favorites, err := h.app.Queries.ListMyFavorites.List(loggedUsrId)
	if err != nil {
		return h.ProperlyError(err)
	}

	h.response.SetStatusCode(http.StatusOK)
	h.response.SetData(favorites)

	return *h.response
}

func (h *Handler) ProperlyError(err error) domain.Response {
	h.response.SetStatusCode(http.StatusInternalServerError)
	h.response.SetMetadata(wrapError(err))

	var notFound *validation.NotFound
	if errors.As(err, &notFound) {
		h.response.SetStatusCode(http.StatusNotFound)
		return *h.response
	}

	var badRequest *validation.BadRequest
	if errors.As(err, &badRequest) {
		h.response.SetStatusCode(http.StatusBadRequest)
		return *h.response
	}

	return *h.response
}

func (h *Handler) validateAuth(r domain.Request) error {
	authorization := r.Headers["Authorization"]
	if authorization == "" {
		authorization = r.Cookies["usess"]
	}

	if authorization == "" {
		return validation.UnauthorizedError("sessão expirada ou inválida")
	}

	_, err := h.jwt.ValidateToken(authorization)
	return err
}

func wrapError(err error) map[string]string {
	return map[string]string{
		"error": err.Error(),
	}
}
