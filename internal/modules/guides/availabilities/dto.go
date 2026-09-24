package availabilities

import (
	"time"

	"github.com/google/uuid"
)

type Availability struct {
	ID        uuid.UUID
	GuideID   string
	StartDate *time.Time
	EndDate   *time.Time
	Reason    *string
}

type CreateAvailabilityRequest struct {
	GuideID   string     `json:"guide_id"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Reason    *string    `json:"reason"`
}
