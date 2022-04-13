package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/handler"
)

func main() {
	lambda.Start(handleRoutes)
}

func handleRoutes(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	h := handler.NewHandler()

	request := domain.ParseRequest(event)

	return domain.WrapResponse(h.Get(request))
}
