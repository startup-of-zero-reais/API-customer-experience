package data

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/expressions"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
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
	var awsConfs []func(*config.LoadOptions) error
	awsConf := append(awsConfs, config.WithRegion("us-east-1"))

	cfg, err := config.LoadDefaultConfig(context.TODO(), awsConf...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	modelDynamo := drivers.NewDynamoClient(
		context.TODO(),
		&domayn.Config{
			TableName: "UserModel",
			Table: table.NewTable(
				"UserModel",
				UserModel{},
			),
			Client: dynamodb.NewFromConfig(cfg),
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

	if len(user) == 0 {
		return nil, validation.NotFoundError("usuário não encontrado")
	}

	return &user[0], nil
}

func (r UserRepository) UpdatePassword(email, password string) error {
	userModel, err := r.Find(email)
	if err != nil {
		return err
	}

	item := r.modelDynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("ID", userModel.ID),
	).AndWhere(
		expressions.NewSortKeyCondition("Email").Equal(email),
	)

	newPassword := providers.NewEncryptProvider()
	passwordCondition := expressions.NewKeyCondition("Password", newPassword.Hash(password))

	item.Update(
		passwordCondition,
	)

	return r.modelDynamo.Perform(drivers.UPDATE, item, &UserModel{})
}
