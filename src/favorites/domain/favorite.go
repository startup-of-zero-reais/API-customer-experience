package domain

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

type (
	// Price represents the price of a meal with value and currency
	Price struct {
		Formatted string `json:"formatted"`
		Value     int    `json:"value"`
	}

	// Meal represents a meal item
	Meal struct {
		ID          string `json:"id"`
		Flavour     string `json:"flavour"`
		Slug        string `json:"-"`
		Ingredients string `json:"ingredients"`
		Photo       string `json:"photo"`
		Price       Price  `json:"price"`
		Self        string `json:"_self"`
		Company     string `json:"-"`
	}

	// Favorite is a struct which represents a favorite meal from a restaurant
	Favorite struct {
		ID      string `json:"id"`
		UserID  string `json:"-"`
		Meal    Meal   `json:"meal"`
		Company string `json:"company"`
	}

	MealRepository interface {
		GetMeal(mealID string) (*Meal, error)
	}

	FavoriteRepository interface {
		Add(favorite *Favorite) error
		UsrFavorites(loggedUsrID string) ([]Favorite, error)
		Delete(loggedUsrID, id string) error
	}
)

func NewFavorite(id, userID, company string, meal *Meal) (*Favorite, error) {
	if id == "" {
		id = uuid.NewString()
	}

	if meal == nil {
		return nil, validation.BadRequestError("o prato favorito deve ser informado")
	}

	if company == "" {
		return nil, validation.BadRequestError("o campo restaurante é obrigatório")
	}

	return &Favorite{
		ID:      id,
		UserID:  userID,
		Meal:    *meal,
		Company: company,
	}, nil
}

func NewMeal(id, flavour, ingredients, photo, self string, price *Price, company string) (*Meal, error) {
	if id == "" {
		return nil, validation.BadRequestError("o campo de identificação é obrigatório")
	}

	if flavour == "" {
		return nil, validation.BadRequestError("o campo sabor é obrigatório")
	}

	if ingredients == "" {
		return nil, validation.BadRequestError("o campo ingredientes é obrigatório")
	}

	if photo == "" {
		return nil, validation.BadRequestError("a foto é obrigatória")
	}

	if price == nil {
		return nil, validation.BadRequestError("o preço é obrigatório")
	}

	return &Meal{
		ID:          id,
		Flavour:     flavour,
		Slug:        slugify(flavour),
		Ingredients: ingredients,
		Photo:       photo,
		Price:       *price,
		Self:        self,
		Company:     company,
	}, nil
}

func NewPrice(value int) (*Price, error) {
	if value == 0 {
		return nil, validation.BadRequestError("o preço é obrigatório")
	}

	brlPrice := number.Decimal(
		float64(value)/float64(100.00),
		number.MinFractionDigits(2),
		number.MaxFractionDigits(2),
	)
	formatted := message.NewPrinter(language.Portuguese).Sprintf("R$ %d", brlPrice)

	return &Price{
		Value:     value,
		Formatted: formatted,
	}, nil
}

func (p *Price) UnmarshalJSON(data []byte) error {
	var value int
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	price, err := NewPrice(value)
	if err != nil {
		return err
	}

	p.Formatted = price.Formatted
	p.Value = price.Value

	return nil
}

func slugify(text string) string {
	return slug.Make(text)
}
