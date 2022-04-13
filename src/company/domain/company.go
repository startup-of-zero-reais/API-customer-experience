package domain

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

type (
	CompanyFood struct {
		ID          string `json:"id"`
		Flavour     string `json:"flavour"`
		Ingredients string `json:"ingredients"`
		Price       int    `json:"price"`
		Photo       string `json:"photo,omitempty"`
		CompanySlug string `json:"_self"`
	}

	ResponsePrice struct {
		Formatted string `json:"formatted"`
		Value     int    `json:"value"`
	}
)

func (c CompanyFood) MarshalJSON() ([]byte, error) {
	return []byte(`{"id":"` + c.ID + `","flavour":"` + c.Flavour + `","ingredients":"` + c.Ingredients + `","price":` + c.getPrice() + `,"photo":"` + c.Photo + `","_self":"` + c.self() + `"}`), nil
}

func (c CompanyFood) self() string {
	return fmt.Sprintf("%s/%s", os.Getenv("APP_URL"), c.CompanySlug)
}

func (c CompanyFood) getPrice() string {
	r := ResponsePrice{}

	num := number.Decimal(float64(c.Price) / float64(100.00))

	r.Formatted = message.NewPrinter(language.Portuguese).Sprintf("R$ %d", num)
	r.Value = c.Price

	result, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(result)

}
