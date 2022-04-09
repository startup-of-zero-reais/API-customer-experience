package handler

import (
	s "github.com/startup-of-zero-reais/API-customer-experience/src/common/service"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/service"
)

type (
	Commands struct {
		SignIn  service.SignIn
		SignOut service.SignOut
	}
	Queries struct {
		GetUser service.GetUser
	}

	Application struct {
		Commands
		Queries
	}
)

func NewApplication(jwtService s.JwtService) *Application {
	usrRepository := data.NewUserRepository()
	sessRepository := data.NewSessionRepository(jwtService)

	return &Application{
		Commands: Commands{
			SignIn: service.NewSignIn(
				usrRepository,
				sessRepository,
			),
			SignOut: service.NewSignOut(
				usrRepository,
				sessRepository,
			),
		},
		Queries: Queries{
			GetUser: service.NewGetUser(usrRepository),
		},
	}
}
