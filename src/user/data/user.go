package data

import (
	"context"
	"errors"

	"github.com/startup-of-zero-reais/dynamo-for-lambda/expressions"

	"github.com/startup-of-zero-reais/API-customer-experience/src/user/domain"
	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/table"
)

type (
	// UserRepository is a interface to access the user data
	UserRepository interface {
		FindByUser(id, email string) (*UserModel, error)
		Save(user *domain.User) error
	}

	// UserRepositoryImpl is a implementation of UserRepository
	UserRepositoryImpl struct {
		eventsDynamo *drivers.DynamoClient
		modelDynamo  *drivers.DynamoClient
	}
)

func NewUserRepository() UserRepository {
	eventsDynamo := drivers.NewDynamoClient(
		context.TODO(),
		&domayn.Config{
			TableName: "UserEvent",
			Table: table.NewTable(
				"UserEvent",
				UserEvent{},
			),
			Endpoint: "http://customer_experience-db:8000",
		},
	)
	eventsDynamo.Migrate()

	modelDynamo := drivers.NewDynamoClient(
		context.TODO(),
		&domayn.Config{
			TableName: "UserModel",
			Table: table.NewTable(
				"UserModel",
				UserModel{},
			),
			Endpoint: "http://customer_experience-db:8000",
		},
	)
	modelDynamo.Migrate()

	return &UserRepositoryImpl{
		eventsDynamo: eventsDynamo,
		modelDynamo:  modelDynamo,
	}
}

func (r *UserRepositoryImpl) FindByUser(id, email string) (*UserModel, error) {
	keyCondition := expressions.NewKeyCondition("ID", id)
	sortKeyCondition := expressions.NewSortKeyCondition("Email").Equal(email)
	sql := r.modelDynamo.NewExpressionBuilder().Where(keyCondition).AndWhere(sortKeyCondition)

	var result UserModel
	err := r.modelDynamo.Perform(drivers.GET, sql, &result)
	if err != nil {
		return &UserModel{}, err
	}

	if result.ID == "" {
		return &UserModel{}, errors.New("usuário não encontrado")
	}

	return &result, nil
}

func (r *UserRepositoryImpl) Save(user *domain.User) error {
	userModel := UserModel{
		ID:        user.ID,
		Name:      user.Name,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  user.Password.Hash(),
		Avatar:    user.Avatar,
		Addresses: "",
	}

	item := r.modelDynamo.NewExpressionBuilder()
	item.SetItem(userModel)

	var result UserModel

	return r.modelDynamo.Put(item, &result)
}
