package agencies

import "time"

type Agency struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Representative string    `json:"representative"`
	Email          string    `json:"email"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateAgencyRequest struct {
	Name           string `json:"name"`
	Representative string `json:"representative"`
	Email          string `json:"email"`
}

type UpdateAgencyRequest struct {
	Name           string `json:"name"`
	Representative string `json:"representative"`
	Email          string `json:"email"`
	IsActive       bool   `json:"is_active"`
}

type ListAgenciesRequest struct {
	Name     *string
	Email    *string
	IsActive *bool
}
