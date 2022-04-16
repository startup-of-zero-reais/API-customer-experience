package data

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/domain"
	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/table"
)

type (
	MealRepositoryImpl struct {
		dynamo *drivers.DynamoClient
	}
)

func NewMealRepository() domain.MealRepository {
	var awsConfs []func(*config.LoadOptions) error
	awsConf := append(awsConfs, config.WithRegion("us-east-1"))

	cfg, err := config.LoadDefaultConfig(context.TODO(), awsConf...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	config := &domayn.Config{
		TableName: "CompanyMeals",
		Table: table.NewTable(
			"CompanyMeals",
			CompanyMeals{},
		),
		Client: dynamodb.NewFromConfig(cfg),
	}

	dynamo := drivers.NewDynamoClient(
		context.TODO(),
		config,
	)
	dynamo.Migrate()

	return &MealRepositoryImpl{
		dynamo: dynamo,
	}
}

func (m *MealRepositoryImpl) GetMeal(mealID string) (*domain.Meal, error) {
	return nil, errors.New("not implemented")
}
