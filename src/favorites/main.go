package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/handler"
)

func main() {
	lambda.Start(handleRoutes)
}

func handleRoutes(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	h := handler.NewHandler()
	request := domain.ParseRequest(event)

	switch event.RequestContext.HTTP.Method {
	case "POST":
		return domain.WrapResponse(h.AddToFavorite(request))
	case "GET":
		return domain.WrapResponse(h.MyFavorites(request))
	case "DELETE":
		return domain.WrapResponse(h.RemoveFavorite(request))
	default:
		panic("Method not implemented")
	}
}
