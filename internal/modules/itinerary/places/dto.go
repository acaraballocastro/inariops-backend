package places

import "time"

type Place struct {
	ID          string
	Name        string
	Description *string
	ZoneID      *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type CreatePlaceRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ZoneID      *string `json:"zone_id"`
	IsActive    bool    `json:"is_active"`
}
