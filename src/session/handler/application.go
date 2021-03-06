package handler

import (
	dt "github.com/startup-of-zero-reais/API-customer-experience/src/common/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	s "github.com/startup-of-zero-reais/API-customer-experience/src/common/service"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/service"
)

type (
	Commands struct {
		SignIn          service.SignIn
		SignOut         service.SignOut
		RecoverPassword service.RecoverPassword
	}
	Queries struct {
		GetUser service.GetUser
	}

	Application struct {
		Commands
		Queries
	}
)

func NewApplication(jwtService s.JwtService, logger *providers.LogProvider) *Application {
	usrRepository := data.NewUserRepository()
	sessRepository := data.NewSessionRepository(jwtService)
	otpRepository := data.NewOTPRepository()
	evtRepository := dt.NewEventRepository()

	sender := data.NewSender()

	return &Application{
		Commands: Commands{
			SignIn: service.NewSignIn(
				usrRepository,
				sessRepository,
				evtRepository,
				logger,
			),
			SignOut: service.NewSignOut(
				usrRepository,
				sessRepository,
				evtRepository,
				logger,
			),
			RecoverPassword: service.NewRecoverPassword(
				usrRepository,
				otpRepository,
				sender,
				evtRepository,
				logger,
			),
		},
		Queries: Queries{
			GetUser: service.NewGetUser(usrRepository, logger),
		},
	}
}
