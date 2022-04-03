package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type (
	// Dynamo is a interface to access the dynamo data
	Dynamo struct {
		ctx    context.Context
		Client *dynamodb.Client
		Tables []string
	}

	Config struct {
		ProvisionedThroughput struct {
			ReadCapacityUnits  int64
			WriteCapacityUnits int64
		}
	}
)

func NewDynamo() *Dynamo {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})

	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	return &Dynamo{
		Client: svc,
		ctx:    ctx,
	}
}

func (d *Dynamo) Init() error {
	tables, err := d.Client.ListTables(d.ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return err
	}

	d.Tables = tables.TableNames

	return nil
}

func (d *Dynamo) CreateTable(name string, config *Config) error {

	out, err := d.Client.CreateTable(d.ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("eventId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("timestamp"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("key"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("eventId"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("timestamp"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("key"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("timestamp"),
						KeyType:       types.KeyTypeHash,
					},
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(config.ProvisionedThroughput.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(config.ProvisionedThroughput.WriteCapacityUnits),
		},
		TableName: aws.String(name),
	})

	if err != nil {
		return err
	}

	d.Tables = append(d.Tables, name)
	fmt.Println(out)

	return nil
}

func (d *Dynamo) Save() error {
	return nil
}
