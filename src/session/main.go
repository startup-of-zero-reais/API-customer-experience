package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/handler"
)

func main() {
	lambda.Start(handleRoutes)
}

func handleRoutes(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handler := handler.NewHandler()

	request := domain.ParseRequest(event)

	signInCase := regexp.MustCompile(`\/sign-in$`).MatchString(event.RequestContext.HTTP.Path)
	signOutCase := regexp.MustCompile(`\/sign-out$`).MatchString(event.RequestContext.HTTP.Path)
	recoverPasswordCase := regexp.MustCompile(`\/recover-password$`).MatchString(event.RequestContext.HTTP.Path)

	switch {
	case signInCase:
		return domain.WrapResponse(handler.SignIn(request))
	case signOutCase:
		return domain.WrapResponse(handler.SignOut(request))
	case recoverPasswordCase:
		return domain.WrapResponse(handler.RecoverPassword(request))
	default:
		log.Println("[ERROR] Invalid path:", event.RequestContext.HTTP.Path)
		r := domain.NewResponse()
		r.SetStatusCode(http.StatusNotFound)
		r.SetMetadata(map[string]string{"error": "endpoint não encontrado"})
		return domain.WrapResponse(*r)
	}

}
