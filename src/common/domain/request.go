package domain

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type (
	// Request represents the parsed request from the API Gateway
	Request struct {
		Headers     map[string]string
		Cookies     map[string]string
		Body        string
		PathParams  map[string]string
		QueryParams map[string]string
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
			log.Println("[ERROR] parsing body with header:", request.Headers["Content-Type"], err)
			panic(err)
		} else {
			request.Body = string(decodedBody)
		}
	}

	if request.Headers["Authorization"] != "" {
		request.Headers["Authorization"] = strings.Replace(request.Headers["Authorization"], "Bearer ", "", 1)
	}

	return Request{
		Headers:     request.Headers,
		Cookies:     cookies,
		Body:        request.Body,
		PathParams:  request.PathParameters,
		QueryParams: request.QueryStringParameters,
	}
}

// WrapResponse returns a response with the status code and the body
func WrapResponse(response Response) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode:      response.StatusCode,
		Headers:         response.Headers,
		Body:            response.Body.ToJson(),
		Cookies:         response.Cookies,
		IsBase64Encoded: false,
	}, nil
}
