package handler

import (
	"encoding/json"
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
		return h.response.HandleError(err)
	}

	loggedUsrId := h.jwt.DecodedToken("id").(string)

	mealInput := struct {
		ID string `json:"meal"`
	}{}

	err = json.Unmarshal([]byte(r.Body), &mealInput)
	if err != nil {
		return h.response.HandleError(err)
	}

	if mealInput.ID == "" {
		return h.response.HandleError(
			validation.BadRequestError("o prato favorito deve ser informado"),
		)
	}

	err = h.app.Commands.AddToFavorite.Meal(loggedUsrId, mealInput.ID)
	if err != nil {
		return h.response.HandleError(err)
	}

	h.response.SetStatusCode(http.StatusCreated)

	return *h.response
}

func (h *Handler) RemoveFavorite(r domain.Request) domain.Response {
	err := h.validateAuth(r)
	if err != nil {
		return h.response.HandleError(err)
	}

	loggedUsrId := h.jwt.DecodedToken("id").(string)

	favorite := struct {
		ID string `json:"favorite"`
	}{
		ID: r.PathParams["favoriteID"],
	}

	if favorite.ID == "" {
		return h.response.HandleError(
			validation.BadRequestError("o favorito deve ser informado"),
		)
	}

	err = h.app.Commands.RemoveFromFavorite.Favorite(loggedUsrId, favorite.ID)
	if err != nil {
		return h.response.HandleError(err)
	}

	h.response.SetStatusCode(http.StatusNoContent)

	return *h.response
}

func (h *Handler) MyFavorites(r domain.Request) domain.Response {
	err := h.validateAuth(r)
	if err != nil {
		return h.response.HandleError(err)
	}

	loggedUsrId := h.jwt.DecodedToken("id").(string)

	favorites, err := h.app.Queries.ListMyFavorites.List(loggedUsrId)
	if err != nil {
		return h.response.HandleError(err)
	}

	h.response.SetStatusCode(http.StatusOK)
	h.response.SetData(favorites)

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
