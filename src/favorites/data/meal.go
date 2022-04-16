package data

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/domain"
	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/expressions"
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

	if os.Getenv("ENVIRONMENT") == "development" {
		awsConf = append(awsConf, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: os.Getenv("ENDPOINT")}, nil
				},
			),
		))
	}

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
	// defer dynamo.FlushDb()
	// err = dynamo.Seed(CompanyMealsSeed()...)
	// if err != nil {
	// 	dynamo.Log.Error("%s\n", err.Error())
	// }

	return &MealRepositoryImpl{
		dynamo: dynamo,
	}
}

func (m *MealRepositoryImpl) GetMeal(mealID string) (*domain.Meal, error) {
	sql := m.dynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("MealID", mealID),
	).SetIndex("MealIndex")

	var companyMeals []CompanyMeals
	err := m.dynamo.Perform(drivers.QUERY, sql, &companyMeals)
	if err != nil {
		return nil, err
	}

	if len(companyMeals) == 0 {
		return nil, validation.NotFoundError("prato não encontrado")
	}

	companyMeal := companyMeals[0]

	mealPrice, err := domain.NewPrice(companyMeal.MealPrice)
	if err != nil {
		return nil, err
	}

	meal, err := domain.NewMeal(
		companyMeal.MealID,
		companyMeal.MealFlavour,
		companyMeal.MealIngredients,
		companyMeal.MealPhoto,
		fmt.Sprintf("%s/%s/%s", os.Getenv("APP_URL"), companyMeal.Slug, companyMeal.MealSlug),
		mealPrice,
		companyMeal.Slug,
	)
	if err != nil {
		return nil, err
	}

	return meal, nil
}
