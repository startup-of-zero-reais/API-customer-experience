package data

import (
	"context"
	"log"
	"os"

	"github.com/startup-of-zero-reais/dynamo-for-lambda/expressions"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/startup-of-zero-reais/API-customer-experience/src/favorites/domain"
	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/table"
)

type (
	FavoritesRepositoryImpl struct {
		dynamo *drivers.DynamoClient
	}
)

func NewFavoritesRepository() domain.FavoriteRepository {
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
		TableName: "Favorites",
		Table: table.NewTable(
			"Favorites",
			Favorites{},
		),
		Client: dynamodb.NewFromConfig(cfg),
	}

	dynamo := drivers.NewDynamoClient(
		context.TODO(),
		config,
	)
	dynamo.Migrate()

	return &FavoritesRepositoryImpl{
		dynamo: dynamo,
	}
}

func (f *FavoritesRepositoryImpl) Add(favorite *domain.Favorite) error {
	fav := Favorites{
		FavoriteID:  favorite.ID,
		UserID:      favorite.UserID,
		CompanySlug: favorite.Meal.Company,
		MealID:      favorite.Meal.ID,
		MealSlug:    favorite.Meal.Slug,
	}

	sql := f.dynamo.NewExpressionBuilder().SetItem(fav)

	err := f.dynamo.Perform(drivers.PUT, sql, &Favorites{})
	if err != nil {
		return err
	}

	return nil
}

func (f *FavoritesRepositoryImpl) UsrFavorites(loggedUsrID string) ([]domain.Favorite, error) {
	sql := f.dynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("UserID", loggedUsrID),
	).SetIndex("MealIndex")

	var dbFavorites []Favorites
	err := f.dynamo.Perform(drivers.QUERY, sql, &dbFavorites)
	if err != nil {
		return nil, err
	}

	favorites := make([]domain.Favorite, len(dbFavorites))
	for i, dbFav := range dbFavorites {
		favorites[i] = domain.Favorite{
			ID:     dbFav.FavoriteID,
			UserID: dbFav.UserID,
			Meal: domain.Meal{
				ID:      dbFav.MealID,
				Company: dbFav.CompanySlug,
				Slug:    dbFav.MealSlug,
			},
		}
	}

	return favorites, nil
}

func (f *FavoritesRepositoryImpl) Delete(loggedUsrID, id string) error {
	sql := f.dynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("UserID", loggedUsrID),
	).AndWhere(
		expressions.NewSortKeyCondition("FavoriteID").Equal(id),
	)

	err := f.dynamo.Perform(drivers.DELETE, sql, &Favorites{})
	if err != nil {
		return err
	}

	return nil
}
