package handler

import (
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
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

func NewApplication(logger *providers.LogProvider) *Application {
	usrRepository := data.NewUserRepository()

	return &Application{
		Commands: Commands{
			CreateUser: service.NewCreateUser(usrRepository, logger),
			UpdateUser: service.NewUpdateUser(usrRepository, logger),
			DeleteUser: service.NewDeleteUser(usrRepository, logger),
		},
		Queries: Queries{
			GetUser: service.NewGetUser(usrRepository, logger),
		},
	}
}
