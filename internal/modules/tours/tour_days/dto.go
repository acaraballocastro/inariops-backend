package tourdays

import (
	"inariops/internal/domain"
	"time"
)

type CreateTourDayInput struct {
	ReservationID          string    `json:"reservation_id" validate:"required"`
	Title                  *string   `json:"title"`
	StartDateTime          time.Time `json:"start_datetime"`
	Duration               *string   `json:"duration"`
	ZoneID                 *string   `json:"zone_id"`
	PeopleCount            *int      `json:"people_count"`
	MeetingPoint           *string   `json:"meeting_point"`
	GuideHotelID           *string   `json:"guide_hotel_id"`
	CustomerHotelID        *string   `json:"customer_hotel_id"`
	GuideLiabilityNotes    *string   `json:"guide_liability_notes"`
	CustomerLiabilityNotes *string   `json:"customer_liability_notes"`
	Remuneration           *float64  `json:"remuneration"`
}

type UpdateTourDayInput struct {
	ID                     string                   `json:"id" validate:"required"`
	ReservationID          string                   `json:"reservation_id" validate:"required"`
	Title                  *string                  `json:"title"`
	StartDateTime          time.Time                `json:"start_datetime"`
	Duration               *string                  `json:"duration"`
	Status                 domain.ReservationStatus `json:"status"`
	ZoneID                 *string                  `json:"zone_id"`
	GuideID                *string                  `json:"guide_id"`
	PeopleCount            *int                     `json:"people_count"`
	MeetingPoint           *string                  `json:"meeting_point"`
	GuideHotelID           *string                  `json:"guide_hotel_id"`
	CustomerHotelID        *string                  `json:"customer_hotel_id"`
	GuideLiabilityNotes    *string                  `json:"guide_liability_notes"`
	CustomerLiabilityNotes *string                  `json:"customer_liability_notes"`
	Remuneration           *float64                 `json:"remuneration"`
	VoucherStatus          domain.VoucherStatus     `json:"voucher_status"`
}

type CancelTourDayInput struct {
	ID string `json:"id" validate:"required"`
}
