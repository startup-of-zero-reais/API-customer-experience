package domain

import (
	"encoding/base64"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type (
	// Request represents the parsed request from the API Gateway
	Request struct {
		Headers map[string]string
		Cookies map[string]string
		Body    string
	}
)

// ParseRequest parses the request from the API Gateway and returns a Request
func ParseRequest(request events.APIGatewayV2HTTPRequest) Request {
	cookies := map[string]string{}

	for _, cookie := range request.Cookies {
		cookieSlices := strings.Split(cookie, "=")
		cookies[cookieSlices[0]] = cookieSlices[1]
	}

	if content := request.Headers["Content-Type"]; content != "application/json" {
		decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			panic(err)
		}

		request.Body = string(decodedBody)
	}

	if request.Headers["Authorization"] != "" {
		request.Headers["Authorization"] = strings.Replace(request.Headers["Authorization"], "Bearer ", "", 1)
	}

	return Request{
		Headers: request.Headers,
		Cookies: cookies,
		Body:    request.Body,
	}
}

// WrapResponse returns a response with the status code and the body
func WrapResponse(response Response) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode:      response.StatusCode,
		Body:            response.Body.ToJson(),
		IsBase64Encoded: false,
	}, nil
}
