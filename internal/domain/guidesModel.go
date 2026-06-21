package domain

import "time"

type Guide struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	MaxToursPerDay int       `json:"max_tours_per_day"`
	CreatedAt      time.Time `json:"created_at"`
}
