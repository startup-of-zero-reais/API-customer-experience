package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/table"
)

type (
	EventRepositoryImpl struct {
		eventsDynamo *drivers.DynamoClient
	}
)

func NewEventRepository() domain.EventRepository {
	var awsConfs []func(*config.LoadOptions) error
	awsConf := append(awsConfs, config.WithRegion("us-east-1"))

	cfg, err := config.LoadDefaultConfig(context.TODO(), awsConf...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	eventsDynamo := drivers.NewDynamoClient(
		context.TODO(),
		&domayn.Config{
			TableName: "UserEvent",
			Table: table.NewTable(
				"UserEvent",
				UserEvent{},
			),
			Client: dynamodb.NewFromConfig(cfg),
		},
	)

	eventsDynamo.Migrate()
	// defer eventsDynamo.FlushDb()

	return &EventRepositoryImpl{
		eventsDynamo: eventsDynamo,
	}
}

func (e *EventRepositoryImpl) Emit(key, eventID string, event domain.Event, data interface{}, timestamp int64) error {
	sql := e.eventsDynamo.NewExpressionBuilder().SetItem(UserEvent{
		Key:       key,
		EventID:   eventID,
		EventType: event,
		Data:      DataToMap(data),
		Timestamp: timestamp,
	})

	return e.eventsDynamo.Perform(drivers.PUT, sql, &UserEvent{})
}

func DataToMap(data interface{}) string {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	mapped := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() == reflect.Ptr && !field.IsZero() {
			// if p, ok := field.Interface().(*fields.Password); ok {
			// 	mapped[v.Type().Field(i).Name] = p.Hash()
			// } else {
			mapped[v.Type().Field(i).Name] = field.Elem().Interface()
			// }
		} else {
			mapped[v.Type().Field(i).Name] = fmt.Sprintf("%v", field.Interface())
		}
	}

	bytes, err := json.Marshal(mapped)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
