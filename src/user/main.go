package main

import (
	"encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/handler"
)

func main() {
	lambda.Start(handleRoutes)
}

func handleRoutes(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	h := handler.NewHandler()

	request := ParseRequestHeaders(event)

	switch event.RequestContext.HTTP.Method {
	case "GET":
		return WrapResponse(h.Get(request.Headers))
	case "POST":
		return WrapResponse(h.Post(request.Body))
	case "PUT":
		return WrapResponse(h.Put())
	case "DELETE":
		return WrapResponse(h.Delete())
	default:
		panic("Method not implemented")
	}
}

func ParseRequestHeaders(request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPRequest {
	if content := request.Headers["Content-Type"]; content != "application/json" {
		decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			panic(err)
		}

		request.Body = string(decodedBody)
	}

	return request
}

func WrapResponse(response domain.Response) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode:      response.StatusCode,
		Body:            response.Body.ToJson(),
		IsBase64Encoded: false,
	}, nil
}
