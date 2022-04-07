package handler

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/data"
	"github.com/startup-of-zero-reais/API-customer-experience/src/user/service"
)

type (
	Commands struct {
		CreateUser service.CreateUser
		UpdateUser service.UpdateUser
		DeleteUser service.DeleteUser
	}

	Queries struct {
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
		Commands: Commands{
			CreateUser: service.NewCreateUser(usrRepository),
			UpdateUser: service.NewUpdateUser(usrRepository),
			DeleteUser: service.NewDeleteUser(usrRepository),
		},
		Queries: Queries{
			GetUser: service.NewGetUser(usrRepository),
		},
	}
}
