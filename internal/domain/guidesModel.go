package domain

import "time"

type Guide struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	MaxToursPerDay int       `json:"max_tours_per_day"`
	CreatedAt      time.Time `json:"created_at"`
}

type GuideUser struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	MaxToursPerDay int       `json:"max_tours_per_day"`
	CreatedAt      time.Time `json:"created_at"`
}
