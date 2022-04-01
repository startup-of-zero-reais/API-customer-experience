package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/handler"
)

func main() {
	lambda.Start(handleRoutes)
}

func handleRoutes(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	h := handler.Handler{}
	switch request.RequestContext.HTTP.Method {
	case "GET":
		return WrapResponse(h.Get())
	case "POST":
		return WrapResponse(h.Post(request.Body))
	case "PUT":
		return WrapResponse(h.Put())
	case "DELETE":
		return WrapResponse(h.Delete())
	default:
		return WrapResponse(h.Get())
	}
}

func WrapResponse(response domain.Response) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode:      response.StatusCode,
		Body:            response.Body.ToJson(),
		IsBase64Encoded: false,
	}, nil
}
