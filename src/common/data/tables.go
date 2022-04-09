package data

import "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"

type (
	UserEvent struct {
		Key       string       `diinamo:"type:string;hash"`
		EventID   string       `diinamo:"type:string;range"`
		EventType domain.Event `diinamo:"type:string;gsi:EventTypeIndex;keyPairs:EventType=EventID"`
		Data      string       `diinamo:"type:string"`
		Timestamp int64        `diinamo:"type:int64;gsi:TimestampIndex;keyPairs:Key=Timestamp"`
	}
)
