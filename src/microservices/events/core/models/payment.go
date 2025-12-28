package models

import "time"

type PaymentEvent struct {
	PaymentID  int       `json:"payment_id" validate:"required"`
	UserID     int       `json:"user_id" validate:"required"`
	Amount     float64   `json:"amount" validate:"required"`
	Status     string    `json:"status" validate:"required"`
	Timestamp  time.Time `json:"timestamp" validate:"required"`
	MethodType *string   `json:"method_type,omitempty"`
}
