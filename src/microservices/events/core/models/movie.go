package models

import "time"

type MovieEvent struct {
	MovieID     int       `json:"movie_id" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Action      string    `json:"action" validate:"required"`
	UserID      *int      `json:"user_id,omitempty"`
	Rating      *float64  `json:"rating,omitempty" validate:"omitempty,min=0,max=10"`
	Genres      []string  `json:"genres,omitempty"`
	Description *string   `json:"description,omitempty"`
	Timestamp   time.Time `json:"timestamp" validate:"required"`
}
