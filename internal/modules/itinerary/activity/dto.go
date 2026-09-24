package activity

import "time"

type Activity struct {
	ID          string
	Name        string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type CreateActivityRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type UpdateActivityRequest struct {
	ID          string  `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}
