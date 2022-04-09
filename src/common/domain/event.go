package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	Event string

	UserEvent struct {
		Key       string      `json:"key"`
		EventID   string      `json:"event_id"`
		Event     Event       `json:"event"`
		Data      interface{} `json:"data"`
		Timestamp int64       `json:"timestamp"`

		repository EventRepository `json:"-"`
	}

	EventEmitter interface {
		Emit(key string, event Event, data interface{}) error
	}

	EventRepository interface {
		Emit(key, eventID string, event Event, data interface{}, timestamp int64) error
	}
)

func NewEventEmitter(repository EventRepository) *UserEvent {
	return &UserEvent{
		repository: repository,
	}
}

func (u *UserEvent) Emit(key string, event Event, data interface{}) error {
	return u.repository.Emit(key, uuid.NewString(), event, data, time.Now().Unix())
}
