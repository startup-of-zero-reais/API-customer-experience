package domain

import (
	"encoding/json"
	"net/http"
)

type (
	Body struct {
		StatusCode int                    `json:"statusCode"`
		Message    string                 `json:"message"`
		Data       interface{}            `json:"data,omitempty"`
		Metadata   map[string]interface{} `json:"metadata,omitempty"`
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

func (r *Response) SetMetadata(metadata map[string]interface{}) *Response {
	r.Body.Data = nil
	r.Body.Metadata = metadata

	return r
}

func (b *Body) ToJson() string {
	jsonBody, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	return string(jsonBody)
}
