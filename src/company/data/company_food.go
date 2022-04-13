package data

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/startup-of-zero-reais/API-customer-experience/src/company/domain"
	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/table"
)

type (
	CompanyFoodsRepository interface {
		FoodList() ([]domain.CompanyFood, error)
	}

	CompanyFoodsRepositoryImpl struct {
		dynamo *drivers.DynamoClient
	}
)

func NewCompanyFoodsRepository() CompanyFoodsRepository {
	var awsConfs []func(*config.LoadOptions) error
	awsConf := append(awsConfs, config.WithRegion("us-east-1"))

	cfg, err := config.LoadDefaultConfig(context.TODO(), awsConf...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	config := &domayn.Config{
		TableName: "CompanyFoodsModel",
		Table: table.NewTable(
			"CompanyFoodsModel",
			CompanyFoodsModel{},
		),
		Client: dynamodb.NewFromConfig(cfg),
	}

	dynamo := drivers.NewDynamoClient(
		context.TODO(),
		config,
	)
	dynamo.Migrate()
	// defer modelDynamo.FlushDb()

	return &CompanyFoodsRepositoryImpl{
		dynamo: dynamo,
	}
}

func (c *CompanyFoodsRepositoryImpl) FoodList() ([]domain.CompanyFood, error) {
	companyFoods := []domain.CompanyFood{
		{
			ID:          "1",
			Flavour:     "Mussarela",
			Ingredients: "Massa, molho de tomate e mussarela",
			Price:       1999,
			Photo:       "https://randomuser.me/api/portraits/lego/1.jpg",
			CompanySlug: "pizzaria-del-vitiente",
		},
		{
			ID:          "2",
			Flavour:     "Calabresa",
			Ingredients: "Massa, molho de tomate, mussarela e calabresa",
			Price:       1999,
			Photo:       "https://randomuser.me/api/portraits/lego/1.jpg",
			CompanySlug: "pizzaria-del-vitiente",
		},
		{
			ID:          "3",
			Flavour:     "Bacon",
			Ingredients: "Massa, molho de tomate, mussarela e bacon",
			Price:       1999,
			Photo:       "https://randomuser.me/api/portraits/lego/1.jpg",
			CompanySlug: "pizzaria-del-vitiente",
		},
	}

	return companyFoods, nil
}
