package handler

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/service"
)

type (
	Commands struct{}
	Queries  struct {
		GetUser service.GetUser
	}

	Application struct {
		Commands
		Queries
	}
)

func NewApplication() *Application {
	usrRepository := data.NewUserRepository()
	return &Application{
		Commands: Commands{},
		Queries: Queries{
			GetUser: service.NewGetUser(usrRepository),
		},
	}
}
