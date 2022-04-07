package main

import (
	"encoding/base64"
	"strings"

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

	request := ParseRequest(event)

	switch event.RequestContext.HTTP.Method {
	case "GET":
		return WrapResponse(h.Get(request))
	case "POST":
		return WrapResponse(h.Post(request))
	case "PUT":
		return WrapResponse(h.Put(request))
	case "DELETE":
		return WrapResponse(h.Delete(request))
	default:
		panic("Method not implemented")
	}
}

func ParseRequest(request events.APIGatewayV2HTTPRequest) handler.Request {
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

	return handler.Request{
		Headers: request.Headers,
		Cookies: cookies,
		Body:    request.Body,
	}
}

func WrapResponse(response domain.Response) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode:      response.StatusCode,
		Body:            response.Body.ToJson(),
		IsBase64Encoded: false,
	}, nil
}
