package domain

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
)

type (
	Body struct {
		StatusCode int         `json:"statusCode"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data,omitempty"`
		Metadata   interface{} `json:"metadata,omitempty"`
	}

	// Response struct represents a response
	Response struct {
		StatusCode        int                 `json:"statusCode"`
		Headers           map[string]string   `json:"headers"`
		MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
		Body              *Body               `json:"body"`
		IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
		Cookies           []string            `json:"cookies"`
	}
)

func NewResponse() *Response {
	return &Response{
		StatusCode:        200,
		Headers:           make(map[string]string),
		MultiValueHeaders: make(map[string][]string),
		Body: &Body{
			StatusCode: 200,
			Message:    "OK",
		},
		IsBase64Encoded: false,
		Cookies:         make([]string, 0),
	}
}

func (r *Response) SetStatusCode(statusCode int) *Response {
	r.StatusCode = statusCode
	r.Body.StatusCode = statusCode
	r.Body.Message = http.StatusText(statusCode)

	return r
}

func (r *Response) SetData(data interface{}) *Response {
	r.Body.Data = data
	r.Body.Metadata = nil

	return r
}

func (r *Response) SetMetadata(metadata interface{}) *Response {
	r.Body.Data = nil
	r.Body.Metadata = metadata

	return r
}

func (r *Response) AddHeader(key, value string) *Response {
	r.Headers[key] = value

	return r
}

func (r *Response) HandleError(err error) Response {
	r.SetStatusCode(http.StatusInternalServerError)
	r.SetMetadata(WrapError(err))

	var notFound *validation.NotFound
	if errors.As(err, &notFound) {
		r.SetStatusCode(http.StatusNotFound)
	}

	var unauthorized *validation.Unauthorized
	if errors.As(err, &unauthorized) {
		r.SetStatusCode(http.StatusUnauthorized)
	}

	var conflict *validation.EntityAlreadyExists
	if errors.As(err, &conflict) {
		r.SetStatusCode(http.StatusConflict)
	}

	var badRequest *validation.BadRequest
	if errors.As(err, &badRequest) {
		r.SetStatusCode(http.StatusBadRequest)
	}

	var fieldValidator *validation.FieldValidator
	if errors.As(err, &fieldValidator) {
		r.SetStatusCode(http.StatusBadRequest)
	}

	return *r
}

func (b *Body) ToJson() string {
	jsonBody, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	return string(jsonBody)
}

func WrapError(err error) map[string]string {
	return map[string]string{
		"error": err.Error(),
	}
}
