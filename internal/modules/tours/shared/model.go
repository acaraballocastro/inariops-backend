package tours

import (
	"inariops/internal/domain"
	"time"
)

type TourDay struct {
	ID            string `json:"id"`
	ReservationID string `json:"reservation_id"`

	Code  *string `json:"code"`
	Title *string `json:"title"`

	StartDateTime time.Time `json:"start_datetime"`
	Duration      *string   `json:"duration"`

	ZoneID  *string `json:"zone_id"`
	GuideID *string `json:"guide_id"`

	Status domain.ReservationStatus `json:"status"`

	PeopleCount *int `json:"people_count"`

	Remuneration *float64 `json:"remuneration"`

	MeetingPoint *string `json:"meeting_point"`

	GuideHotelID    *string `json:"guide_hotel_id"`
	CustomerHotelID *string `json:"customer_hotel_id"`

	GuideLiabilityNotes    *string `json:"guide_liability_notes"`
	CustomerLiabilityNotes *string `json:"customer_liability_notes"`

	VoucherStatus domain.VoucherStatus `json:"voucher_status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
