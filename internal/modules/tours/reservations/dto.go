package reservations

import (
	"inariops/internal/domain"
	"time"
)

type Reservation struct {
	ID               string                   `json:"id"`
	Code             *string                  `json:"code"`
	Title            *string                  `json:"title"`
	Description      *string                  `json:"description"`
	Status           domain.ReservationStatus `json:"status"`
	QuoteNumber      *string                  `json:"quote_number"`
	FileNumber       *string                  `json:"file_number"`
	AgencyID         *string                  `json:"agency_id"`
	TotalPeopleCount *int                     `json:"total_people_count"`
	StartDate        time.Time                `json:"start_date"`
	EndDate          time.Time                `json:"end_date"`
	SignatureStatus  domain.SignatureStatus   `json:"signature_status"`
	VoucherStatus    domain.VoucherStatus     `json:"voucher_status"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
}

type CreateReservationRequest struct {
	Title            *string   `json:"title"`
	Description      *string   `json:"description"`
	AgencyID         *string   `json:"agency_id"`
	TotalPeopleCount *int      `json:"total_people_count"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
}

type UpdateReservationRequest struct {
	Code             string
	Title            *string   `json:"title"`
	Description      *string   `json:"description"`
	QuoteNumber      *string   `json:"quote_number"`
	FileNumber       *string   `json:"file_number"`
	AgencyID         *string   `json:"agency_id"`
	TotalPeopleCount *int      `json:"total_people_count"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
