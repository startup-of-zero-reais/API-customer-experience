package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/handler"
)

func main() {
	lambda.Start(handleRoutes)
}

func handleRoutes(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	l := providers.NewLogProvider()
	l.LoggerConfig(event)
	h := handler.NewHandler(l)
	request := domain.ParseRequest(event)

	responseHandler := handleResponseWithLogger(l)

	switch event.RequestContext.HTTP.Method {
	case "POST":
		return responseHandler(h.AddToFavorite(request))
	case "GET":
		return responseHandler(h.MyFavorites(request))
	case "DELETE":
		return responseHandler(h.RemoveFavorite(request))
	default:
		panic("Method not implemented")
	}
}

// handleResponseWithLogger is a helper function to handle the response and log the response
func handleResponseWithLogger(logger *providers.LogProvider) func(response domain.Response) (events.APIGatewayV2HTTPResponse, error) {
	return func(response domain.Response) (events.APIGatewayV2HTTPResponse, error) {
		logger.LogResponse(response)

		return domain.WrapResponse(response)
	}
}
