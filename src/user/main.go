package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/handler"
)

func main() {
	lambda.Start(handleRoutes)
}

func handleRoutes(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	logger := providers.NewLogProvider()
	logger.LoggerConfig(event)
	h := handler.NewHandler(logger)

	request := domain.ParseRequest(event)

	switch event.RequestContext.HTTP.Method {
	case "GET":
		return domain.WrapResponse(h.Get(request))
	case "POST":
		return domain.WrapResponse(h.Post(request))
	case "PUT":
		return domain.WrapResponse(h.Put(request))
	case "DELETE":
		return domain.WrapResponse(h.Delete(request))
	default:
		panic("Method not implemented")
	}
}
