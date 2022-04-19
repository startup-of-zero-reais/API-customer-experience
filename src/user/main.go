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

	request, err := domain.ParseRequest(event)
	if err != nil {
		logger.Errorln("[ERROR] parsing request:", err)
		panic(err)
	}

	responseHandler := handleResponseWithLogger(logger)

	switch event.RequestContext.HTTP.Method {
	case "GET":
		return responseHandler(h.Get(request))
	case "POST":
		return responseHandler(h.Post(request))
	case "PUT":
		return responseHandler(h.Put(request))
	case "DELETE":
		return responseHandler(h.Delete(request))
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
