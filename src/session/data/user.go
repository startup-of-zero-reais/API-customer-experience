package data

import (
	"context"

	"github.com/startup-of-zero-reais/dynamo-for-lambda/expressions"

	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"
	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/table"
)

type (
	UserRepository struct {
		modelDynamo *drivers.DynamoClient
	}
)

func NewUserRepository() UserRepository {
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

	return UserRepository{
		modelDynamo: modelDynamo,
	}
}

func (u UserRepository) Find(email string) (*domain.User, error) {
	sql := u.modelDynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("Email", email),
	).SetIndex("EmailIndex")

	var user []domain.User
	err := u.modelDynamo.Perform(drivers.QUERY, sql, &user)
	if err != nil {
		return nil, err
	}

	return &user[0], nil
}
