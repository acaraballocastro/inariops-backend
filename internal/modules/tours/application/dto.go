package reservationsapp

import (
	"inariops/internal/modules/customers"
	"inariops/internal/modules/tours/reservations"
	tours "inariops/internal/modules/tours/shared"
	"time"
)

type ReservationDetail struct {
	Reservation reservations.Reservation `json:"reservation"`
	Customers   []customers.Customer     `json:"customers"`
	TourDays    []tours.TourDay          `json:"tour_days"`
}

type CreateReservationRequest struct {
	Title            *string                      `json:"title"`
	Description      *string                      `json:"description"`
	AgencyID         *string                      `json:"agency_id"`
	TotalPeopleCount *int                         `json:"total_people_count"`
	StartDate        time.Time                    `json:"start_date"`
	EndDate          time.Time                    `json:"end_date"`
	Customers        []CustomerReservationRequest `json:"customers"`
}

type CustomerReservationRequest struct {
	ID             *string `json:"id,omitempty"`
	FullName       string  `json:"full_name"`
	DocumentNumber string  `json:"document_number"`
	Phone          string  `json:"phone"`
	Email          string  `json:"email"`
	Age            int     `json:"age"`
}

type UpdateReservationRequest struct {
	Code             string                       `json:"code"`
	Title            *string                      `json:"title"`
	Description      *string                      `json:"description"`
	AgencyID         *string                      `json:"agency_id"`
	TotalPeopleCount *int                         `json:"total_people_count"`
	QuoteNumber      *string                      `json:"quote_number"`
	FileNumber       *string                      `json:"file_number"`
	StartDate        time.Time                    `json:"start_date"`
	EndDate          time.Time                    `json:"end_date"`
	Customers        []CustomerReservationRequest `json:"customers"`
}
