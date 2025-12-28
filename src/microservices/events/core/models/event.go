package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	UserCreated    string = "users.created"
	MovieCreated   string = "movies.created"
	PaymentCreated string = "payment.created"
)

type Event[T any] struct {
	EventID   string    `json:"event_id" binding:"required"`
	EventType string    `json:"event_type" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
	Payload   T         `json:"payload" binding:"required"`
}

func NewEvent[T any](eventType string, payload T) *Event[T] {
	return &Event[T]{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

func (e *Event[T]) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event[T]) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}
